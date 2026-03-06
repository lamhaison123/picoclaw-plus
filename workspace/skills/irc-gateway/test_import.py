#!/usr/bin/env python3
"""Quick test to check if gateway can be imported"""

import sys
import os

print("=" * 60)
print("Testing IRC Gateway Import")
print("=" * 60)

# Test 1: Check Python version
print(f"\n1. Python version: {sys.version}")

# Test 2: Check telegram library
try:
    import telegram
    print(f"2. telegram library: ✅ v{telegram.__version__}")
except ImportError as e:
    print(f"2. telegram library: ❌ {e}")
    print("   Run: pip3 install python-telegram-bot")
    sys.exit(1)

# Test 3: Check config file
config_path = os.path.expanduser("~/.picoclaw/config.json")
if os.path.exists(config_path):
    print(f"3. PicoClaw config: ✅ {config_path}")
    
    # Try to read token
    import json
    try:
        with open(config_path) as f:
            config = json.load(f)
            if "channels" in config and "telegram" in config["channels"]:
                token = config["channels"]["telegram"].get("token", "")
                if token:
                    print(f"   Token found: {token[:10]}...{token[-10:]}")
                else:
                    print("   ⚠️  Token is empty")
            else:
                print("   ⚠️  Telegram not configured")
    except Exception as e:
        print(f"   ⚠️  Error reading config: {e}")
else:
    print(f"3. PicoClaw config: ❌ Not found at {config_path}")

# Test 4: Try to import gateway module (without running main)
print("\n4. Testing gateway.py import...")
try:
    # Read gateway.py and check syntax
    with open("gateway.py") as f:
        code = f.read()
    
    # Try to compile it
    compile(code, "gateway.py", "exec")
    print("   ✅ Gateway syntax is valid")
    
    # Try to execute without running main
    code_no_main = code.replace('if __name__ == "__main__":', 'if False:')
    exec(compile(code_no_main, "gateway.py", "exec"))
    print("   ✅ Gateway imports successfully")
    
except SyntaxError as e:
    print(f"   ❌ Syntax error: {e}")
    sys.exit(1)
except Exception as e:
    print(f"   ❌ Import error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

print("\n" + "=" * 60)
print("✅ All checks passed! Gateway is ready to run.")
print("=" * 60)
print("\nTo start gateway:")
print("  python3 gateway.py")
print()
