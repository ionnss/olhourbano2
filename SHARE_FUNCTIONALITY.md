# Share Functionality for Report Detail Page

## Overview

The share functionality allows users to share reports as images on social media platforms. When a user clicks the "Compartilhar" (Share) button on a report detail page, a modal opens with various sharing options.

## Features

### Share Button
- Located next to the "Ver no Mapa" (View on Map) button
- Styled with a green gradient background
- Responsive design for mobile devices

### Share Modal
- Displays sharing options for different social media platforms
- Shows a preview of the generated share image
- Includes options for:
  - WhatsApp
  - Facebook
  - Twitter/X
  - Telegram
  - Copy Link

### Share Image Generation
- Uses html2canvas library to generate social media-friendly images
- Creates 1200x630px images (optimal for social media)
- Includes:
  - Report category and icon
  - Report description (truncated if too long)
  - Location information
  - Vote count
  - Creation date
  - Olho Urbano branding

## Technical Implementation

### Files Modified/Created

1. **Templates**
   - `templates/components/04_report_detail_content.html` - Added share button
   - `templates/pages/04_report_detail.html` - Added html2canvas library and share.js

2. **CSS**
   - `static/css/report.css` - Added share button and modal styles

3. **JavaScript**
   - `static/js/share.js` - Complete share functionality implementation

4. **Backend**
   - `handlers/api.go` - Added ShareImageHandler (placeholder for future server-side generation)
   - `routes/routes.go` - Added share image API route

### Dependencies

- **html2canvas**: Client-side image generation library
- **Bootstrap Icons**: For social media icons
- **Google Fonts**: Montserrat font for image generation

### Browser Compatibility

The share functionality works in modern browsers that support:
- ES6+ JavaScript features
- Canvas API
- Fetch API
- CSS Grid and Flexbox

## Usage

1. Navigate to any report detail page
2. Click the "Compartilhar" button
3. Choose a sharing option from the modal
4. The selected platform will open with pre-filled content

## Customization

### Image Design
The share image design can be customized by modifying the `generateShareImage()` function in `share.js`. The current design includes:

- Gradient background with subtle pattern
- Category icon and name
- Report description and location
- Vote count display
- Olho Urbano branding

### Social Media Platforms
Additional platforms can be added by:
1. Adding a new share option in the modal HTML
2. Creating a corresponding share function
3. Adding CSS styles for the new platform

### API Integration
For server-side image generation, implement the `ShareImageHandler` in `handlers/api.go` using a library like:
- `github.com/fogleman/gg` (Go)
- `github.com/golang/freetype` (Go)

## Error Handling

The implementation includes error handling for:
- Missing html2canvas library
- Failed image generation
- Network errors during sharing
- Unsupported browser features

## Future Enhancements

1. **Server-side Image Generation**: Implement proper server-side image generation for better performance
2. **Image Caching**: Cache generated images to avoid regeneration
3. **Custom Templates**: Allow different image templates for different report types
4. **Analytics**: Track sharing metrics
5. **QR Codes**: Add QR codes to share images for easy mobile access

## Testing

To test the functionality:
1. Start the application
2. Navigate to a report detail page
3. Click the share button
4. Verify that the modal opens and image preview is generated
5. Test each sharing option
6. Test on mobile devices for responsiveness
