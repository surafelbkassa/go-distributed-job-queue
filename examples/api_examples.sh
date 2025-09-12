#!/bin/bash

# Example script to demonstrate the job queue API

BASE_URL="http://localhost:8080/api/v1"

echo "=== Job Queue API Examples ==="
echo ""

# Health check
echo "1. Health Check:"
curl -s "${BASE_URL}/health" | jq '.'
echo -e "\n"

# Get job types
echo "2. Available Job Types:"
curl -s "${BASE_URL}/job-types" | jq '.'
echo -e "\n"

# Create email job
echo "3. Creating Email Job:"
EMAIL_JOB=$(curl -s -X POST "${BASE_URL}/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {
      "to": "user@example.com",
      "subject": "Welcome to our service!",
      "body": "Thank you for signing up!"
    },
    "priority": 1
  }')
echo "$EMAIL_JOB" | jq '.'
EMAIL_JOB_ID=$(echo "$EMAIL_JOB" | jq -r '.job_id')
echo -e "\n"

# Create image resize job
echo "4. Creating Image Resize Job:"
IMAGE_JOB=$(curl -s -X POST "${BASE_URL}/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "image_resize",
    "payload": {
      "image_url": "https://example.com/image.jpg",
      "width": 800,
      "height": 600
    }
  }')
echo "$IMAGE_JOB" | jq '.'
IMAGE_JOB_ID=$(echo "$IMAGE_JOB" | jq -r '.job_id')
echo -e "\n"

# Create report generation job
echo "5. Creating Report Generation Job:"
REPORT_JOB=$(curl -s -X POST "${BASE_URL}/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "generate_report",
    "payload": {
      "type": "sales",
      "period": "monthly",
      "year": 2023,
      "month": 12
    }
  }')
echo "$REPORT_JOB" | jq '.'
REPORT_JOB_ID=$(echo "$REPORT_JOB" | jq -r '.job_id')
echo -e "\n"

# Wait a bit for jobs to process
echo "6. Waiting 3 seconds for jobs to process..."
sleep 3
echo ""

# Check job status
echo "7. Checking Email Job Status:"
curl -s "${BASE_URL}/jobs/${EMAIL_JOB_ID}" | jq '.'
echo -e "\n"

# List all jobs
echo "8. Listing All Jobs:"
curl -s "${BASE_URL}/jobs?limit=10" | jq '.'
echo -e "\n"

# List pending jobs
echo "9. Listing Pending Jobs:"
curl -s "${BASE_URL}/jobs?status=pending&limit=5" | jq '.'
echo -e "\n"

# List completed jobs
echo "10. Listing Completed Jobs:"
curl -s "${BASE_URL}/jobs?status=completed&limit=5" | jq '.'
echo -e "\n"

# Get statistics
echo "11. Queue Statistics:"
curl -s "${BASE_URL}/stats" | jq '.'
echo -e "\n"

echo "=== Example completed! ==="