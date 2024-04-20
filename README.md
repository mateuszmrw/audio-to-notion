# Audio to Notion

This project automates the transcription of audio files into text, analyzes the content using OpenAI's API, and then creates a new page in a Notion database with the analyzed data. The system leverages AWS Lambda for processing and uses S3 for audio file storage.

## Architecture Overview

- **Whisper**: Transcribes user-provided audio (.MP3) files to text.
- **OpenAI API**: Processes the transcript to generate structured JSON data.
- **Notion API**: Uses the JSON to create a new page in a Notion database.
- **AWS Lambda**: Hosts the processing logic.
- **AWS S3 Bucket**: Stores the audio files and triggers the Lambda function.
- **Terraform**: Manages the AWS infrastructure setup.

## Prerequisites

- AWS Account
- Notion Account
- OpenAI API Key

## Environment Variables

Before deploying the project, ensure you have set up the following environment variables:

- `OPEN_AI_KEY`: Your OpenAI API key.
- `NOTION_TOKEN`: Integration token from Notion.
- `NOTION_PAGE_ID`: ID of the Notion database where pages will be added.
- `AWS_REGION`: Your AWS region.

## Project Structure

- `/shared`: Contains shared Golang code for the Lambda function.
- `/lambda`: Main entry point for the Lambda function.
- `/infra`: Terraform scripts for setting up AWS infrastructure including Lambda functions and S3 buckets.

## Setup Instructions

### Step 1: Prepare your AWS environment

1. Navigate to the `/infra` directory.
2. Use Terraform to deploy the infrastructure:
   ```bash
   terraform init
   terraform apply
Confirm the setup when prompted.

### Step 2: Configure your Lambda function

Navigate to the `/lambda` folder.
Ensure the Lambda function is set up with the necessary environment variables.

### Step 3: Upload an audio file

Upload an .MP3 audio file to the designated S3 bucket. This will trigger the Lambda function to start the transcription and data processing workflow.

### Step 4: Check Notion

Once the Lambda function has processed the audio file, check your Notion database for a new page with the analyzed data.
