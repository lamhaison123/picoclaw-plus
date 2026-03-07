#!/usr/bin/env python3
"""
IRC-Style Communication Gateway for PicoClaw
Integrates with PicoClaw's Team system for mention-based routing

This gateway extends PicoClaw's existing Telegram channel to add IRC-style
mention routing (@architect, @developer, etc.) that maps to team roles.
"""

import os
import re
import json
import asyncio
import logging
import subprocess
from typing import Dict, List, Set, Optional
from datetime import datetime
from pathlib import Path

try:
    from telegram import Update
    from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes
except ImportError:
    print("Installing python-telegram-bot...")
    os.system("pip install python-telegram-bot")
    from telegram import Update
    from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes

# Configuration
REPO_ROOT = Path(__file__).parent.parent.parent.parent
ENV_FILE = REPO_ROOT / ".env"
CONFIG_FILE = Path.home() / ".picoclaw" / "config.json"

def load_config():
    """Load configuration from PicoClaw config.json or .env"""
    bot_token = ""
    picoclaw_bin = "picoclaw"
    team_id = "irc-dev-team"
    
    # Try loading from PicoClaw config first
    if CONFIG_FILE.exists():
        try:
            with open(CONFIG_FILE) as f:
                config = json.load(f)
                # Get Telegram token from channels.telegram.token
                if "channels" in config and "telegram" in config["channels"]:
                    bot_token = config["channels"]["telegram"].get("token", "")
                    logger.info(f"✅ Loaded Telegram token from {CONFIG_FILE}")
        except Exception as e:
            logger.warning(f"Failed to load config.json: {e}")
    
    # Fallback to .env file
    if not bot_token and ENV_FILE.exists():
        with open(ENV_FILE) as f:
            for line in f:
                if line.strip() and not line.startswith('#'):
                    if '=' in line:
                        key, value = line.strip().split('=', 1)
                        if key == "TELEGRAM_BOT_TOKEN":
                            bot_token = value.strip('"').strip("'")
                            logger.info(f"✅ Loaded Telegram token from {ENV_FILE}")
                        elif key == "PICOCLAW_BIN":
                            picoclaw_bin = value.strip('"').strip("'")
                        elif key == "IRC_TEAM_ID":
                            team_id = value.strip('"').strip("'")
    
    # Environment variables override
    bot_token = os.getenv("TELEGRAM_BOT_TOKEN", bot_token)
    picoclaw_bin = os.getenv("PICOCLAW_BIN", picoclaw_bin)
    team_id = os.getenv("IRC_TEAM_ID", team_id)
    
    return bot_token, picoclaw_bin, team_id

BOT_TOKEN, PICOCLAW_BIN, TEAM_ID = load_config()

# Agent/Role configuration (maps mentions to team roles)
ROLE_MAP = {
    "architect": {"emoji": "🏗️", "description": "System design & architecture"},
    "developer": {"emoji": "💻", "description": "Code implementation"},
    "tester": {"emoji": "🧪", "description": "Testing & QA"},
    "manager": {"emoji": "📋", "description": "Project coordination"}
}

# Session management
sessions: Dict[int, Dict] = {}
agent_status: Dict[str, str] = {role: "idle" for role in ROLE_MAP}

# Logging
logging.basicConfig(
    format='[%(asctime)s] %(levelname)s: %(message)s',
    level=logging.INFO,
    datefmt='%H:%M:%S'
)
logger = logging.getLogger(__name__)


def load_team_config() -> Optional[Dict]:
    """Load team configuration from PicoClaw"""
    team_config_path = Path.home() / ".picoclaw" / "workspace" / "teams" / f"{TEAM_ID}.json"
    if team_config_path.exists():
        with open(team_config_path) as f:
            return json.load(f)
    return None


def extract_mentions(text: str) -> Set[str]:
    """Extract @mentions from message text"""
    mentions = set()
    for role in ROLE_MAP:
        if f"@{role}" in text.lower():
            mentions.add(role)
    return mentions


def generate_session_id(chat_id: int) -> str:
    """Generate short session ID for IRC-style formatting"""
    timestamp = datetime.now().strftime("%H%M%S")
    return f"{abs(hash(chat_id)) % 10000:04d}{timestamp[-4:]}"


async def execute_with_team(role: str, message: str, session_id: str, chat_id: str) -> str:
    """Execute task using PicoClaw team system with specific role"""
    try:
        agent_status[role] = "busy"
        logger.info(f"[{session_id}] Routing to @{role}: {message[:50]}...")
        
        # Use picoclaw team execute command with role-based routing
        # The team system will automatically route to the correct agent
        cmd = [
            PICOCLAW_BIN, "team", "execute", TEAM_ID,
            "-t", message,
            "--role", role,  # Specify which role should handle this
            "--format", "json"
        ]
        
        process = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )
        
        stdout, stderr = await asyncio.wait_for(process.communicate(), timeout=120)
        
        if process.returncode == 0:
            try:
                result = json.loads(stdout.decode())
                response = result.get("result", stdout.decode().strip())
            except json.JSONDecodeError:
                response = stdout.decode().strip()
        else:
            error_msg = stderr.decode().strip()
            logger.error(f"Team execution failed: {error_msg}")
            response = f"Error: {error_msg}"
        
        agent_status[role] = "idle"
        return response
        
    except asyncio.TimeoutError:
        agent_status[role] = "idle"
        return "⏱️ Task timeout (120s exceeded)"
    except FileNotFoundError:
        agent_status[role] = "idle"
        return f"❌ Error: picoclaw binary not found at '{PICOCLAW_BIN}'"
    except Exception as e:
        agent_status[role] = "idle"
        logger.exception(f"Error executing with team: {e}")
        return f"❌ Error: {str(e)}"


async def format_response(role: str, response: str, session_id: str) -> str:
    """Format response in IRC style"""
    emoji = ROLE_MAP[role]["emoji"]
    return f"[{session_id}] {emoji} {role.upper()}: {response}"


async def start_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /start command"""
    chat_id = update.effective_chat.id
    session_id = generate_session_id(chat_id)
    
    sessions[chat_id] = {
        "session_id": session_id,
        "history": []
    }
    
    # Check if team exists
    team_status = "✅ Connected" if check_team_exists() else "⚠️ Not configured"
    
    welcome = f"""🤖 **PicoClaw IRC Gateway** [{session_id}]

**Team:** {TEAM_ID} ({team_status})

**Available roles:**
{chr(10).join(f"  @{name} {info['emoji']} - {info['description']}" for name, info in ROLE_MAP.items())}

**Usage:**
• Tag roles: `@architect design a REST API`
• Multiple: `@architect @developer implement auth`
• Default (no tag): Routes to @manager

**Commands:**
• /who - List active roles
• /status - Check role status
• /team - Show team info
• /clear - Clear session history

**Integration:**
This gateway routes mentions to PicoClaw's Team system.
Configure your team at: `~/.picoclaw/workspace/teams/{TEAM_ID}.json`
"""
    await update.message.reply_text(welcome, parse_mode='Markdown')


def check_team_exists() -> bool:
    """Check if the team is configured in PicoClaw"""
    try:
        result = subprocess.run(
            [PICOCLAW_BIN, "team", "status", TEAM_ID],
            capture_output=True,
            timeout=5
        )
        return result.returncode == 0
    except:
        return False


async def who_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /who command"""
    role_list = "\n".join(
        f"  {info['emoji']} @{name} - {info['description']}"
        for name, info in ROLE_MAP.items()
    )
    await update.message.reply_text(f"**Active Roles:**\n{role_list}", parse_mode='Markdown')


async def status_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /status command"""
    status_list = "\n".join(
        f"  @{role}: {'🔴 busy' if status == 'busy' else '🟢 idle'}"
        for role, status in agent_status.items()
    )
    await update.message.reply_text(f"**Role Status:**\n{status_list}", parse_mode='Markdown')


async def team_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /team command - show team configuration"""
    try:
        result = subprocess.run(
            [PICOCLAW_BIN, "team", "status", TEAM_ID],
            capture_output=True,
            text=True,
            timeout=5
        )
        
        if result.returncode == 0:
            info = f"**Team Status:**\n```\n{result.stdout}\n```"
        else:
            info = f"⚠️ Team '{TEAM_ID}' not found.\n\nCreate it with:\n`picoclaw team create templates/teams/development-team.json`"
        
        await update.message.reply_text(info, parse_mode='Markdown')
    except Exception as e:
        await update.message.reply_text(f"❌ Error: {str(e)}")


async def clear_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle /clear command"""
    chat_id = update.effective_chat.id
    if chat_id in sessions:
        sessions[chat_id]["history"] = []
        await update.message.reply_text("✅ Session history cleared")
    else:
        await update.message.reply_text("⚠️ No active session")


async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle incoming messages with mention-based routing"""
    chat_id = update.effective_chat.id
    message_text = update.message.text
    
    # Initialize session if needed
    if chat_id not in sessions:
        session_id = generate_session_id(chat_id)
        sessions[chat_id] = {"session_id": session_id, "history": []}
    
    session = sessions[chat_id]
    session_id = session["session_id"]
    
    # Extract mentions
    mentions = extract_mentions(message_text)
    
    # Default to manager if no mentions
    if not mentions:
        mentions = {"manager"}
    
    # Store in history
    session["history"].append({
        "user": message_text,
        "roles": list(mentions),
        "timestamp": datetime.now()
    })
    
    # Process roles in parallel
    tasks = []
    for role in mentions:
        tasks.append(execute_with_team(role, message_text, session_id, str(chat_id)))
    
    # Wait for all roles to respond
    responses = await asyncio.gather(*tasks)
    
    # Format and send responses
    for role, response in zip(mentions, responses):
        formatted = await format_response(role, response, session_id)
        await update.message.reply_text(formatted)


async def error_handler(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Handle errors"""
    logger.error(f"Update {update} caused error {context.error}")


def main():
    """Main entry point"""
    if not BOT_TOKEN:
        print("❌ Error: TELEGRAM_BOT_TOKEN not found in .env file")
        print(f"📝 Please add TELEGRAM_BOT_TOKEN=your_token to: {ENV_FILE}")
        return
    
    logger.info("🚀 Starting PicoClaw IRC Gateway...")
    logger.info(f"📁 Repo root: {REPO_ROOT}")
    logger.info(f"🤖 Team ID: {TEAM_ID}")
    logger.info(f"👥 Roles: {', '.join(ROLE_MAP.keys())}")
    
    # Check if picoclaw is accessible
    try:
        result = subprocess.run([PICOCLAW_BIN, "version"], capture_output=True, timeout=5)
        if result.returncode == 0:
            logger.info(f"✅ PicoClaw binary found: {PICOCLAW_BIN}")
        else:
            logger.warning(f"⚠️ PicoClaw binary may not be working: {PICOCLAW_BIN}")
    except:
        logger.error(f"❌ Cannot execute picoclaw binary: {PICOCLAW_BIN}")
        logger.error("   Make sure picoclaw is installed and in PATH")
        return
    
    # Create application
    app = Application.builder().token(BOT_TOKEN).build()
    
    # Register handlers
    app.add_handler(CommandHandler("start", start_command))
    app.add_handler(CommandHandler("who", who_command))
    app.add_handler(CommandHandler("status", status_command))
    app.add_handler(CommandHandler("team", team_command))
    app.add_handler(CommandHandler("clear", clear_command))
    app.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_message))
    app.add_error_handler(error_handler)
    
    # Start bot
    logger.info("✅ Gateway is running. Press Ctrl+C to stop.")
    app.run_polling(allowed_updates=Update.ALL_TYPES)


if __name__ == "__main__":
    main()
