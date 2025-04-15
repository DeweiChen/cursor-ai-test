#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080/api"

# Create a chatroom
echo "Creating a chatroom..."
CHATROOM_RESPONSE=$(curl -s -X POST "$BASE_URL/chatrooms" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Chatroom",
    "description": "A test chatroom for message testing"
  }')

# Extract chatroom ID from response
CHATROOM_ID=$(echo $CHATROOM_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Created chatroom with ID: $CHATROOM_ID"

# Post a message to the chatroom
echo -e "\nPosting a message to the chatroom..."
curl -X POST "$BASE_URL/chatrooms/$CHATROOM_ID/messages" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "TestUser",
    "content": "Hello, this is a test message!"
  }'

# Post another message
echo -e "\nPosting another message..."
curl -X POST "$BASE_URL/chatrooms/$CHATROOM_ID/messages" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "AnotherUser",
    "content": "Hi there! This is another test message."
  }'

# Get all messages from the chatroom
echo -e "\nGetting all messages from the chatroom..."
curl -s "$BASE_URL/chatrooms/$CHATROOM_ID/messages" | jq 