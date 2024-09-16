# File Sharing Platform
## Overview
This is a file-sharing platform built using Go, PostgreSQL, and local file storage. The platform enables users to securely upload, manage, and share files. The system uses JWT-based authentication to ensure that users can only access their own files. The project is designed with efficiency and scalability in mind, leveraging Goroutines for concurrency and handling large file uploads.

## Features
**1. User Authentication:**
</br>JWT-based authentication for secure login and registration.
</br>Endpoints for user registration and login.
</br>Users can manage their own files only.

**2. File Upload and Management:**
</br>Upload files using multipart form.
</br>Metadata for each file (name, size, upload date) is stored in PostgreSQL.
</br>Public URL generation for file access.
</br>Concurrent file uploads using Goroutines.

**3. File Retrieval and Sharing:**
</br>Retrieve uploaded file metadata and generate sharable public URLs.

**4. File Search:**
</br>Search files by name, date, or type with efficient query handling.

## API Endpoints
**Authentication** -
        </br>POST /register: Register a new user.
        </br>POST /login: Log in and receive a JWT token.</br>
</br>**File Operations** -
        </br>POST /upload: Upload a file (requires JWT).
        </br>GET /files: Retrieve all uploaded file metadata (requires JWT).
        </br>GET /share/: Generate a public URL for file sharing (requires JWT).
