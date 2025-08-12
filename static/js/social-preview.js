/**
 * Social Media Preview Testing Utility
 * Helps validate Open Graph and Twitter Card meta tags
 */

class SocialPreviewTester {
    constructor() {
        this.platforms = {
            facebook: 'https://developers.facebook.com/tools/debug/',
            twitter: 'https://cards-dev.twitter.com/validator',
            linkedin: 'https://www.linkedin.com/post-inspector/',
            whatsapp: 'https://wa.me/',
            telegram: 'https://t.me/'
        };
    }

    /**
     * Get current page meta tags for validation
     */
    getMetaTags() {
        const metaTags = {
            title: document.title,
            description: this.getMetaContent('description'),
            ogTitle: this.getMetaContent('og:title'),
            ogDescription: this.getMetaContent('og:description'),
            ogImage: this.getMetaContent('og:image'),
            ogUrl: this.getMetaContent('og:url'),
            twitterTitle: this.getMetaContent('twitter:title'),
            twitterDescription: this.getMetaContent('twitter:description'),
            twitterImage: this.getMetaContent('twitter:image'),
            twitterCard: this.getMetaContent('twitter:card')
        };

        return metaTags;
    }

    /**
     * Get meta content by property/name
     */
    getMetaContent(property) {
        const meta = document.querySelector(`meta[property="${property}"], meta[name="${property}"]`);
        return meta ? meta.getAttribute('content') : null;
    }

    /**
     * Validate meta tags and show results
     */
    validateMetaTags() {
        const tags = this.getMetaTags();
        const issues = [];
        const warnings = [];

        // Check required tags
        if (!tags.ogTitle) issues.push('Missing og:title');
        if (!tags.ogDescription) issues.push('Missing og:description');
        if (!tags.ogImage) issues.push('Missing og:image');
        if (!tags.ogUrl) issues.push('Missing og:url');

        // Check Twitter tags
        if (!tags.twitterTitle) warnings.push('Missing twitter:title (will fallback to og:title)');
        if (!tags.twitterDescription) warnings.push('Missing twitter:description (will fallback to og:description)');
        if (!tags.twitterImage) warnings.push('Missing twitter:image (will fallback to og:image)');

        // Check description length
        if (tags.ogDescription && tags.ogDescription.length > 160) {
            warnings.push(`Description too long (${tags.ogDescription.length} chars, max 160)`);
        }

        // Check image URL format
        if (tags.ogImage) {
            this.validateImageUrl(tags.ogImage, warnings);
            this.checkImageDimensions(tags.ogImage, warnings);
        }

        return { issues, warnings, tags };
    }

    /**
     * Check image dimensions for optimal social sharing
     */
    async checkImageDimensions(imageUrl, warnings) {
        try {
            const img = new Image();
            img.onload = function() {
                if (this.width < 1200 || this.height < 630) {
                    warnings.push(`Image dimensions (${this.width}x${this.height}) are below recommended 1200x630`);
                }
                if (this.width > 1200 || this.height > 630) {
                    warnings.push(`Image dimensions (${this.width}x${this.height}) are larger than recommended 1200x630`);
                }
                if (this.width / this.height !== 1200 / 630) {
                    warnings.push(`Image aspect ratio (${(this.width / this.height).toFixed(2)}) differs from recommended 1.91:1`);
                }
            };
            img.onerror = function() {
                warnings.push(`Could not load image: ${imageUrl}`);
            };
            img.src = imageUrl;
        } catch (error) {
            warnings.push('Could not check image dimensions');
        }
    }

    /**
     * Validate image URL format
     */
    validateImageUrl(imageUrl, warnings) {
        if (!imageUrl) {
            warnings.push('No image URL provided');
            return;
        }
        
        // Check if URL is absolute
        if (!imageUrl.startsWith('http://') && !imageUrl.startsWith('https://')) {
            warnings.push('Image URL should be absolute (start with http:// or https://)');
        }
        
        // Check file extension
        const extension = imageUrl.split('.').pop()?.toLowerCase();
        if (extension && !['jpg', 'jpeg', 'png', 'webp'].includes(extension)) {
            warnings.push(`Image format .${extension} may not be supported by all social platforms. Use JPG, PNG, or WebP.`);
        }
        
        // Check for SVG (not well supported)
        if (extension === 'svg') {
            warnings.push('SVG format is not well supported by social media platforms. Use PNG or JPG instead.');
        }
    }

    /**
     * Open platform testing tools
     */
    openTestingTool(platform) {
        const currentUrl = encodeURIComponent(window.location.href);
        const toolUrl = this.platforms[platform];
        
        if (platform === 'whatsapp') {
            // For WhatsApp, we'll create a test message
            const message = encodeURIComponent(`Testando preview: ${document.title}`);
            window.open(`${toolUrl}?text=${message}%20${currentUrl}`);
        } else if (platform === 'telegram') {
            // For Telegram, we'll create a test message
            const message = encodeURIComponent(`Testando preview: ${document.title}`);
            window.open(`${toolUrl}?text=${message}%20${currentUrl}`);
        } else {
            // For other platforms, open their testing tools
            window.open(`${toolUrl}?q=${currentUrl}`);
        }
    }

    /**
     * Generate preview HTML for testing
     */
    generatePreviewHTML() {
        const tags = this.getMetaTags();
        
        return `
            <div class="social-preview-card">
                <div class="preview-image">
                    <img src="${tags.ogImage || '/static/resource/full_logo.svg'}" alt="Preview image">
                </div>
                <div class="preview-content">
                    <div class="preview-url">${tags.ogUrl || window.location.href}</div>
                    <div class="preview-title">${tags.ogTitle || tags.title}</div>
                    <div class="preview-description">${tags.ogDescription || tags.description}</div>
                </div>
            </div>
        `;
    }

    /**
     * Show validation results in console
     */
    logValidationResults() {
        const results = this.validateMetaTags();
        
        console.group('🔍 Social Media Preview Validation');
        
        if (results.issues.length === 0 && results.warnings.length === 0) {
            console.log('✅ All meta tags are properly configured!');
        }
        
        if (results.issues.length > 0) {
            console.error('❌ Issues found:', results.issues);
        }
        
        if (results.warnings.length > 0) {
            console.warn('⚠️ Warnings:', results.warnings);
        }
        
        console.log('📋 Current meta tags:', results.tags);
        console.groupEnd();
        
        return results;
    }
}

// Initialize the tester
const socialPreviewTester = new SocialPreviewTester();

// Auto-validate on page load (only in development)
if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    document.addEventListener('DOMContentLoaded', () => {
        setTimeout(() => {
            socialPreviewTester.logValidationResults();
        }, 1000);
    });
}

// Make it available globally for manual testing
window.socialPreviewTester = socialPreviewTester;
