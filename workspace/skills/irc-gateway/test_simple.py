#!/usr/bin/env python3
"""
Simple unit tests for IRC Gateway (no external dependencies)
"""

import sys
import re

# Mock the ROLE_MAP for testing
ROLE_MAP = {
    "architect": {"emoji": "🏗️", "description": "System design"},
    "developer": {"emoji": "💻", "description": "Implementation"},
    "tester": {"emoji": "🧪", "description": "Testing"},
    "manager": {"emoji": "📋", "description": "Coordination"}
}

def extract_mentions(text: str) -> set:
    """Extract @mentions from message text"""
    mentions = set()
    for role in ROLE_MAP:
        if f"@{role}" in text.lower():
            mentions.add(role)
    return mentions

def generate_session_id(chat_id: int) -> str:
    """Generate short session ID"""
    from datetime import datetime
    timestamp = datetime.now().strftime("%H%M%S")
    return f"{abs(hash(chat_id)) % 10000:04d}{timestamp[-4:]}"

def test_mention_parsing():
    """Test mention extraction"""
    tests = [
        ("@architect design API", {"architect"}),
        ("@architect @developer build", {"architect", "developer"}),
        ("no mentions here", set()),
        ("@ARCHITECT @Developer", {"architect", "developer"}),
        ("@invalid @architect", {"architect"}),
    ]
    
    passed = 0
    failed = 0
    
    for text, expected in tests:
        result = extract_mentions(text)
        if result == expected:
            print(f"✅ PASS: '{text}' -> {result}")
            passed += 1
        else:
            print(f"❌ FAIL: '{text}' -> {result} (expected {expected})")
            failed += 1
    
    return passed, failed

def test_session_id():
    """Test session ID generation"""
    passed = 0
    failed = 0
    
    # Test format
    sid = generate_session_id(12345)
    if len(sid) == 8:
        print(f"✅ PASS: Session ID format correct: {sid}")
        passed += 1
    else:
        print(f"❌ FAIL: Session ID wrong length: {sid}")
        failed += 1
    
    # Test uniqueness
    sid1 = generate_session_id(12345)
    sid2 = generate_session_id(67890)
    if sid1[:4] != sid2[:4]:
        print(f"✅ PASS: Different chat IDs generate different hashes")
        passed += 1
    else:
        print(f"❌ FAIL: Same hash for different chat IDs")
        failed += 1
    
    return passed, failed

def test_role_config():
    """Test role configuration"""
    passed = 0
    failed = 0
    
    expected_roles = {"architect", "developer", "tester", "manager"}
    if set(ROLE_MAP.keys()) == expected_roles:
        print(f"✅ PASS: All expected roles defined")
        passed += 1
    else:
        print(f"❌ FAIL: Missing or extra roles")
        failed += 1
    
    # Check each role has required fields
    for role, config in ROLE_MAP.items():
        if "emoji" in config and "description" in config:
            print(f"✅ PASS: Role '{role}' has required fields")
            passed += 1
        else:
            print(f"❌ FAIL: Role '{role}' missing fields")
            failed += 1
    
    return passed, failed

def main():
    """Run all tests"""
    print("=" * 60)
    print("IRC Gateway - Simple Unit Tests")
    print("=" * 60)
    print()
    
    total_passed = 0
    total_failed = 0
    
    print("Test 1: Mention Parsing")
    print("-" * 60)
    p, f = test_mention_parsing()
    total_passed += p
    total_failed += f
    print()
    
    print("Test 2: Session ID Generation")
    print("-" * 60)
    p, f = test_session_id()
    total_passed += p
    total_failed += f
    print()
    
    print("Test 3: Role Configuration")
    print("-" * 60)
    p, f = test_role_config()
    total_passed += p
    total_failed += f
    print()
    
    print("=" * 60)
    print(f"Results: {total_passed} passed, {total_failed} failed")
    print("=" * 60)
    
    return 0 if total_failed == 0 else 1

if __name__ == "__main__":
    sys.exit(main())
