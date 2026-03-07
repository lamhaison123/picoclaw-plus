#!/usr/bin/env python3
"""
Test suite for IRC Gateway
Tests mention parsing, role routing, and formatting
"""

import unittest
import asyncio
from unittest.mock import Mock, patch, AsyncMock
from gateway import (
    extract_mentions,
    generate_session_id,
    format_response,
    ROLE_MAP
)


class TestMentionParsing(unittest.TestCase):
    """Test mention extraction from messages"""
    
    def test_single_mention(self):
        """Test extracting single @mention"""
        text = "@architect design a REST API"
        mentions = extract_mentions(text)
        self.assertEqual(mentions, {"architect"})
    
    def test_multiple_mentions(self):
        """Test extracting multiple @mentions"""
        text = "hey @architect @developer let's build this"
        mentions = extract_mentions(text)
        self.assertEqual(mentions, {"architect", "developer"})
    
    def test_no_mentions(self):
        """Test message with no mentions"""
        text = "what's the status?"
        mentions = extract_mentions(text)
        self.assertEqual(mentions, set())
    
    def test_case_insensitive(self):
        """Test mentions are case-insensitive"""
        text = "@ARCHITECT @Developer @TeSTeR"
        mentions = extract_mentions(text)
        self.assertEqual(mentions, {"architect", "developer", "tester"})
    
    def test_invalid_mentions(self):
        """Test invalid mentions are ignored"""
        text = "@architect @invalid @developer"
        mentions = extract_mentions(text)
        self.assertEqual(mentions, {"architect", "developer"})
    
    def test_mention_in_middle(self):
        """Test mentions in middle of text"""
        text = "can @architect help with @developer on this?"
        mentions = extract_mentions(text)
        self.assertEqual(mentions, {"architect", "developer"})


class TestSessionManagement(unittest.TestCase):
    """Test session ID generation"""
    
    def test_session_id_format(self):
        """Test session ID has correct format"""
        chat_id = 12345
        session_id = generate_session_id(chat_id)
        self.assertIsInstance(session_id, str)
        self.assertEqual(len(session_id), 8)
    
    def test_session_id_unique(self):
        """Test different chat IDs generate different session IDs"""
        id1 = generate_session_id(12345)
        id2 = generate_session_id(67890)
        self.assertNotEqual(id1, id2)
    
    def test_session_id_consistent(self):
        """Test same chat ID generates consistent hash part"""
        id1 = generate_session_id(12345)
        id2 = generate_session_id(12345)
        # First 4 chars (hash) should be same
        self.assertEqual(id1[:4], id2[:4])


class TestResponseFormatting(unittest.IsolatedAsyncioTestCase):
    """Test IRC-style response formatting"""
    
    async def test_format_response(self):
        """Test response formatting with role and emoji"""
        response = await format_response("architect", "Here's the design", "test1234")
        self.assertIn("[test1234]", response)
        self.assertIn("🏗️", response)
        self.assertIn("ARCHITECT:", response)
        self.assertIn("Here's the design", response)
    
    async def test_format_all_roles(self):
        """Test formatting works for all defined roles"""
        for role in ROLE_MAP:
            response = await format_response(role, "test message", "sess0001")
            self.assertIn(f"[sess0001]", response)
            self.assertIn(ROLE_MAP[role]["emoji"], response)
            self.assertIn(role.upper(), response)


class TestRoleConfiguration(unittest.TestCase):
    """Test role configuration"""
    
    def test_all_roles_defined(self):
        """Test all expected roles are defined"""
        expected_roles = {"architect", "developer", "tester", "manager"}
        self.assertEqual(set(ROLE_MAP.keys()), expected_roles)
    
    def test_role_has_emoji(self):
        """Test each role has an emoji"""
        for role, config in ROLE_MAP.items():
            self.assertIn("emoji", config)
            self.assertIsInstance(config["emoji"], str)
            self.assertGreater(len(config["emoji"]), 0)
    
    def test_role_has_description(self):
        """Test each role has a description"""
        for role, config in ROLE_MAP.items():
            self.assertIn("description", config)
            self.assertIsInstance(config["description"], str)
            self.assertGreater(len(config["description"]), 0)


class TestIntegration(unittest.IsolatedAsyncioTestCase):
    """Integration tests with mocked PicoClaw"""
    
    @patch('gateway.subprocess.run')
    async def test_check_team_exists(self, mock_run):
        """Test team existence check"""
        from gateway import check_team_exists
        
        # Mock successful team check
        mock_run.return_value = Mock(returncode=0)
        self.assertTrue(check_team_exists())
        
        # Mock failed team check
        mock_run.return_value = Mock(returncode=1)
        self.assertFalse(check_team_exists())
    
    @patch('gateway.asyncio.create_subprocess_exec')
    async def test_execute_with_team_success(self, mock_subprocess):
        """Test successful team execution"""
        from gateway import execute_with_team
        
        # Mock successful execution
        mock_process = AsyncMock()
        mock_process.communicate.return_value = (
            b'{"result": "Design completed"}',
            b''
        )
        mock_process.returncode = 0
        mock_subprocess.return_value = mock_process
        
        result = await execute_with_team("architect", "design API", "test123", "12345")
        self.assertIn("Design completed", result)
    
    @patch('gateway.asyncio.create_subprocess_exec')
    async def test_execute_with_team_timeout(self, mock_subprocess):
        """Test team execution timeout"""
        from gateway import execute_with_team
        
        # Mock timeout
        mock_process = AsyncMock()
        mock_process.communicate.side_effect = asyncio.TimeoutError()
        mock_subprocess.return_value = mock_process
        
        result = await execute_with_team("architect", "design API", "test123", "12345")
        self.assertIn("timeout", result.lower())
    
    @patch('gateway.asyncio.create_subprocess_exec')
    async def test_execute_with_team_error(self, mock_subprocess):
        """Test team execution error handling"""
        from gateway import execute_with_team
        
        # Mock error
        mock_process = AsyncMock()
        mock_process.communicate.return_value = (
            b'',
            b'Team not found'
        )
        mock_process.returncode = 1
        mock_subprocess.return_value = mock_process
        
        result = await execute_with_team("architect", "design API", "test123", "12345")
        self.assertIn("Error", result)


def run_tests():
    """Run all tests"""
    loader = unittest.TestLoader()
    suite = unittest.TestSuite()
    
    # Add all test classes
    suite.addTests(loader.loadTestsFromTestCase(TestMentionParsing))
    suite.addTests(loader.loadTestsFromTestCase(TestSessionManagement))
    suite.addTests(loader.loadTestsFromTestCase(TestResponseFormatting))
    suite.addTests(loader.loadTestsFromTestCase(TestRoleConfiguration))
    suite.addTests(loader.loadTestsFromTestCase(TestIntegration))
    
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)
    
    return result.wasSuccessful()


if __name__ == "__main__":
    success = run_tests()
    exit(0 if success else 1)
