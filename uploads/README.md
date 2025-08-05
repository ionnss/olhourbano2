# Uploads Directory

This directory stores user-uploaded files from report submissions.

## Security Notes

- All uploaded files have their metadata cleaned for security
- Files are renamed with SHA-256 hashes to prevent conflicts
- Only specific file types are allowed (images, PDFs, documents)
- File size limits are enforced server-side

## Git Ignore

This directory is configured in `.gitignore` to:
- ✅ Ignore all uploaded files (`uploads/*`)
- ✅ Keep directory structure (`.gitkeep`)
- ✅ Keep this documentation (`README.md`)

## Production Deployment

In production environments:
- Ensure proper file permissions (readable by web server)
- Consider using a CDN or separate file storage service
- Regular backups of uploaded files
- Monitor disk space usage

## File Structure

Files are organized by upload date and hashed names:
```
uploads/
├── .gitkeep
├── README.md
└── [hashed-filename].[extension]
```