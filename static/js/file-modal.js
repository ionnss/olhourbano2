// File Modal functionality
let currentFileIndex = 0;
let currentFiles = [];

// Image zoom variables
let currentZoom = 1;
let minZoom = 0.1;
let maxZoom = 5;
let isDragging = false;
let lastMouseX = 0;
let lastMouseY = 0;
let translateX = 0;
let translateY = 0;
let initialDistance = 0;
let initialZoom = 1;
let isZooming = false;

// Fallback modal function if Bootstrap is not available
function showModalFallback() {
    const modal = document.getElementById('fileModal');
    if (modal) {
        modal.style.display = 'block';
        modal.classList.add('show');
        document.body.classList.add('modal-open');
        
        // Add backdrop
        const backdrop = document.createElement('div');
        backdrop.className = 'modal-backdrop fade show';
        backdrop.id = 'modalBackdrop';
        document.body.appendChild(backdrop);
    }
}

function hideModalFallback() {
    const modal = document.getElementById('fileModal');
    const backdrop = document.getElementById('modalBackdrop');
    
    if (modal) {
        modal.style.display = 'none';
        modal.classList.remove('show');
        document.body.classList.remove('modal-open');
    }
    
    if (backdrop) {
        backdrop.remove();
    }
}

// Open file modal for single file
function openFileModal(filePath, fileName, fileType) {
    console.log('openFileModal called:', { filePath, fileName, fileType });
    
    const modal = document.getElementById('fileModal');
    const modalTitle = document.getElementById('fileModalLabel');
    
    console.log('Modal elements found:', {
        modal: !!modal,
        modalTitle: !!modalTitle
    });
    
    // Hide all content viewers
    hideAllViewers();
    
    // Show appropriate viewer based on file type
    if (fileType === 'image' || isImageFile(fileName)) {
        showImageViewer(filePath);
    } else if (fileType === 'video' || isVideoFile(fileName)) {
        showVideoViewer(filePath);
    } else if (fileType === 'pdf' || isPdfFile(fileName)) {
        showPdfViewer(filePath);
    } else {
        showDocumentViewer(filePath, fileName);
    }
    
    // Show modal
    if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
        const bsModal = new bootstrap.Modal(modal);
        bsModal.show();
    } else {
        showModalFallback();
    }
}

// Open file gallery for multiple files
function openFileGallery(element) {
    const photos = element.getAttribute('data-photos').split(',');
    const reportId = element.getAttribute('data-report-id');
    
    currentFiles = photos.map(photo => ({
        path: '/' + photo.trim(),
        name: photo.trim(),
        type: getFileType(photo.trim())
    }));
    currentFileIndex = 0;
    
    if (currentFiles.length > 0) {
        openFileModal(currentFiles[0].path, currentFiles[0].name, currentFiles[0].type);
    }
}

// Show image viewer with zoom functionality
function showImageViewer(filePath) {
    const imageViewer = document.getElementById('imageViewer');
    const modalImage = document.getElementById('modalImage');
    const imageLoading = document.getElementById('imageLoading');
    const zoomInstructions = document.getElementById('zoomInstructions');
    
    // Reset zoom state
    resetImageZoom();
    
    // Show loading state
    if (imageLoading) {
        imageLoading.style.display = 'block';
    }
    
    // Set image source
    modalImage.src = filePath;
    imageViewer.style.display = 'block';
    
    // Show zoom instructions on mobile
    if (zoomInstructions && isMobileDevice()) {
        zoomInstructions.classList.remove('hidden');
        // Hide instructions after 5 seconds
        setTimeout(() => {
            zoomInstructions.classList.add('hidden');
        }, 5000);
    }
    
    // Initialize zoom functionality after image loads
    modalImage.onload = function() {
        // Hide loading state
        if (imageLoading) {
            imageLoading.style.display = 'none';
        }
        initializeImageZoom();
    };
    
    // Handle image load error
    modalImage.onerror = function() {
        console.error('Failed to load image:', filePath);
        // Hide loading state
        if (imageLoading) {
            imageLoading.style.display = 'none';
        }
        showImageError();
    };
}

// Show image error state
function showImageError() {
    const imageViewer = document.getElementById('imageViewer');
    const modalImage = document.getElementById('modalImage');
    
    if (imageViewer && modalImage) {
        imageViewer.innerHTML = `
            <div class="image-error text-center p-4">
                <i class="bi bi-exclamation-triangle display-1 text-warning mb-3"></i>
                <h5 class="text-muted">Erro ao carregar imagem</h5>
                <p class="text-muted">Não foi possível carregar esta imagem.</p>
            </div>
        `;
    }
}

// Initialize image zoom functionality
function initializeImageZoom() {
    const modalImage = document.getElementById('modalImage');
    const imageContainer = document.getElementById('imageContainer');
    const zoomInBtn = document.getElementById('zoomIn');
    const zoomOutBtn = document.getElementById('zoomOut');
    const resetZoomBtn = document.getElementById('resetZoom');
    const fitToScreenBtn = document.getElementById('fitToScreen');
    const zoomLevelDisplay = document.getElementById('zoomLevel');
    
    if (!modalImage || !imageContainer) return;
    
    // Reset zoom state
    resetImageZoom();
    
    // Remove existing event listeners to prevent duplicates
    removeZoomEventListeners();
    
    // Zoom control event listeners
    if (zoomInBtn) {
        zoomInBtn.addEventListener('click', () => zoomImage(1.2));
    }
    
    if (zoomOutBtn) {
        zoomOutBtn.addEventListener('click', () => zoomImage(0.8));
    }
    
    if (resetZoomBtn) {
        resetZoomBtn.addEventListener('click', resetImageZoom);
    }
    
    if (fitToScreenBtn) {
        fitToScreenBtn.addEventListener('click', fitImageToScreen);
    }
    
    // Mouse wheel zoom with better handling
    imageContainer.addEventListener('wheel', handleWheelZoom, { passive: false });
    
    // Mouse drag panning with improved handling
    modalImage.addEventListener('mousedown', startDragging);
    document.addEventListener('mousemove', handleDragging);
    document.addEventListener('mouseup', stopDragging);
    
    // Touch events for mobile with better support
    modalImage.addEventListener('touchstart', handleTouchStart, { passive: false });
    modalImage.addEventListener('touchmove', handleTouchMove, { passive: false });
    modalImage.addEventListener('touchend', handleTouchEnd);
    
    // Keyboard shortcuts
    document.addEventListener('keydown', handleKeyboardZoom);
    
    // Update zoom level display
    updateZoomLevelDisplay();
    
    // Initial fit to screen with delay to ensure image is loaded
    setTimeout(fitImageToScreen, 200);
    
    // Add double-click to reset zoom
    modalImage.addEventListener('dblclick', resetImageZoom);
    
    // Add loading class for smooth transitions
    modalImage.classList.add('loading');
    setTimeout(() => {
        modalImage.classList.remove('loading');
    }, 300);
}

// Remove zoom event listeners to prevent duplicates
function removeZoomEventListeners() {
    const modalImage = document.getElementById('modalImage');
    const imageContainer = document.getElementById('imageContainer');
    
    if (modalImage) {
        modalImage.removeEventListener('mousedown', startDragging);
        modalImage.removeEventListener('touchstart', handleTouchStart);
        modalImage.removeEventListener('touchmove', handleTouchMove);
        modalImage.removeEventListener('touchend', handleTouchEnd);
        modalImage.removeEventListener('dblclick', resetImageZoom);
    }
    
    if (imageContainer) {
        imageContainer.removeEventListener('wheel', handleWheelZoom);
    }
}

// Handle mouse wheel zoom with improved precision
function handleWheelZoom(e) {
    e.preventDefault();
    
    const delta = e.deltaY > 0 ? 0.9 : 1.1;
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    zoomImageAtPoint(delta, x, y);
}

// Handle keyboard zoom with better key detection
function handleKeyboardZoom(e) {
    if (document.getElementById('imageViewer').style.display === 'none') return;
    
    // Check if user is typing in an input field
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
    
    switch(e.key) {
        case '+':
        case '=':
            e.preventDefault();
            zoomImage(1.2);
            break;
        case '-':
        case '_':
            e.preventDefault();
            zoomImage(0.8);
            break;
        case '0':
            e.preventDefault();
            resetImageZoom();
            break;
        case 'f':
        case 'F':
            e.preventDefault();
            fitImageToScreen();
            break;
        case 'Escape':
            e.preventDefault();
            closeModal();
            break;
        case 'ArrowLeft':
            e.preventDefault();
            previousFile();
            break;
        case 'ArrowRight':
            e.preventDefault();
            nextFile();
            break;
    }
}

// Close modal function
function closeModal() {
    if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
        const modal = bootstrap.Modal.getInstance(document.getElementById('fileModal'));
        if (modal) {
            modal.hide();
        }
    } else {
        hideModalFallback();
    }
}

// Start dragging with improved handling
function startDragging(e) {
    if (currentZoom <= 1) return;
    
    isDragging = true;
    lastMouseX = e.clientX || e.touches[0].clientX;
    lastMouseY = e.clientY || e.touches[0].clientY;
    e.preventDefault();
    
    const modalImage = document.getElementById('modalImage');
    if (modalImage) {
        modalImage.style.cursor = 'grabbing';
    }
}

// Handle dragging with improved performance
function handleDragging(e) {
    if (!isDragging) return;
    
    const clientX = e.clientX || (e.touches && e.touches[0] ? e.touches[0].clientX : lastMouseX);
    const clientY = e.clientY || (e.touches && e.touches[0] ? e.touches[0].clientY : lastMouseY);
    
    const deltaX = clientX - lastMouseX;
    const deltaY = clientY - lastMouseY;
    
    translateX += deltaX;
    translateY += deltaY;
    
    lastMouseX = clientX;
    lastMouseY = clientY;
    
    updateImageTransform();
}

// Stop dragging with cleanup
function stopDragging() {
    isDragging = false;
    
    const modalImage = document.getElementById('modalImage');
    if (modalImage) {
        modalImage.style.cursor = 'grab';
    }
}

// Touch start handler with improved detection
function handleTouchStart(e) {
    if (e.touches.length === 1) {
        // Single touch - start dragging
        isDragging = true;
        lastMouseX = e.touches[0].clientX;
        lastMouseY = e.touches[0].clientY;
    } else if (e.touches.length === 2) {
        // Two touches - start pinch zoom
        const touch1 = e.touches[0];
        const touch2 = e.touches[1];
        initialDistance = Math.sqrt(
            Math.pow(touch2.clientX - touch1.clientX, 2) +
            Math.pow(touch2.clientY - touch1.clientY, 2)
        );
        initialZoom = currentZoom;
        isZooming = true;
    }
}

// Touch move handler with improved performance
function handleTouchMove(e) {
    e.preventDefault();
    
    if (e.touches.length === 1 && isDragging) {
        // Single touch dragging
        const deltaX = e.touches[0].clientX - lastMouseX;
        const deltaY = e.touches[0].clientY - lastMouseY;
        
        translateX += deltaX;
        translateY += deltaY;
        
        lastMouseX = e.touches[0].clientX;
        lastMouseY = e.touches[0].clientY;
        
        updateImageTransform();
    } else if (e.touches.length === 2 && isZooming) {
        // Two touches - pinch zoom
        const touch1 = e.touches[0];
        const touch2 = e.touches[1];
        const currentDistance = Math.sqrt(
            Math.pow(touch2.clientX - touch1.clientX, 2) +
            Math.pow(touch2.clientY - touch1.clientY, 2)
        );
        
        if (initialDistance > 0) {
            const scale = currentDistance / initialDistance;
            const newZoom = Math.max(minZoom, Math.min(maxZoom, initialZoom * scale));
            zoomImage(newZoom / currentZoom);
        }
    }
}

// Touch end handler with cleanup
function handleTouchEnd(e) {
    isDragging = false;
    isZooming = false;
    initialDistance = 0;
    
    const modalImage = document.getElementById('modalImage');
    if (modalImage) {
        modalImage.style.cursor = 'grab';
    }
}

// Zoom image at specific point with improved precision
function zoomImageAtPoint(scale, x, y) {
    const newZoom = Math.max(minZoom, Math.min(maxZoom, currentZoom * scale));
    const zoomRatio = newZoom / currentZoom;
    
    // Calculate new position to zoom towards the mouse point
    translateX = x - (x - translateX) * zoomRatio;
    translateY = y - (y - translateY) * zoomRatio;
    
    currentZoom = newZoom;
    updateImageTransform();
    updateZoomLevelDisplay();
}

// Zoom image with improved bounds checking
function zoomImage(scale) {
    const newZoom = Math.max(minZoom, Math.min(maxZoom, currentZoom * scale));
    currentZoom = newZoom;
    updateImageTransform();
    updateZoomLevelDisplay();
}

// Reset image zoom with improved state management
function resetImageZoom() {
    currentZoom = 1;
    translateX = 0;
    translateY = 0;
    isDragging = false;
    isZooming = false;
    updateImageTransform();
    updateZoomLevelDisplay();
}

// Fit image to screen with improved calculation
function fitImageToScreen() {
    const modalImage = document.getElementById('modalImage');
    const imageContainer = document.getElementById('imageContainer');
    
    if (!modalImage || !imageContainer) return;
    
    const containerRect = imageContainer.getBoundingClientRect();
    const imageRect = modalImage.getBoundingClientRect();
    
    // Wait for image to be properly loaded
    if (imageRect.width === 0 || imageRect.height === 0) {
        setTimeout(fitImageToScreen, 100);
        return;
    }
    
    const scaleX = containerRect.width / imageRect.width;
    const scaleY = containerRect.height / imageRect.height;
    const scale = Math.min(scaleX, scaleY, 1); // Don't scale up beyond 100%
    
    currentZoom = scale;
    translateX = 0;
    translateY = 0;
    updateImageTransform();
    updateZoomLevelDisplay();
}

// Update image transform with improved performance
function updateImageTransform() {
    const modalImage = document.getElementById('modalImage');
    if (!modalImage) return;
    
    // Apply transform with better performance
    modalImage.style.transform = `translate(${translateX}px, ${translateY}px) scale(${currentZoom})`;
    
    // Add zooming class for animation
    if (!modalImage.classList.contains('zooming')) {
        modalImage.classList.add('zooming');
        setTimeout(() => {
            modalImage.classList.remove('zooming');
        }, 200);
    }
}

// Update zoom level display with improved formatting
function updateZoomLevelDisplay() {
    const zoomLevelDisplay = document.getElementById('zoomLevel');
    if (zoomLevelDisplay) {
        const percentage = Math.round(currentZoom * 100);
        zoomLevelDisplay.textContent = `${percentage}%`;
        
        // Add visual feedback for zoom level
        zoomLevelDisplay.className = 'zoom-level';
        if (percentage < 100) {
            zoomLevelDisplay.classList.add('zoomed-out');
        } else if (percentage > 100) {
            zoomLevelDisplay.classList.add('zoomed-in');
        }
        
        // Update ARIA live region
        zoomLevelDisplay.setAttribute('aria-label', `Zoom: ${percentage}%`);
    }
}

// Show video viewer
function showVideoViewer(filePath) {
    const videoViewer = document.getElementById('videoViewer');
    const modalVideo = document.getElementById('modalVideo');
    const videoSource = document.getElementById('videoSource');
    
    videoSource.src = filePath;
    videoSource.type = getVideoMimeType(filePath);
    modalVideo.load();
    videoViewer.style.display = 'block';
}

// Show PDF viewer
function showPdfViewer(filePath) {
    console.log('showPdfViewer called with:', filePath);
    
    const pdfViewer = document.getElementById('pdfViewer');
    const pdfFrame = document.getElementById('pdfFrame');
    const pdfFallback = document.querySelector('.pdf-fallback');
    const pdfLoadingIndicator = document.getElementById('pdfLoadingIndicator');
    
    console.log('PDF elements found:', {
        pdfViewer: !!pdfViewer,
        pdfFrame: !!pdfFrame,
        pdfFallback: !!pdfFallback,
        pdfLoadingIndicator: !!pdfLoadingIndicator
    });
    
    // Check if mobile device
    const isMobile = isMobileDevice();
    console.log('Is mobile device:', isMobile);
    
    if (isMobile) {
        // On mobile, immediately show option to open in new tab
        console.log('Mobile device detected - showing PDF open option');
        showPdfFallback();
        
        // Set up the open link
        const pdfOpenLink = document.getElementById('pdfOpenLink');
        if (pdfOpenLink) pdfOpenLink.href = filePath;
        
        // Hide the iframe and loading indicator
        if (pdfFrame) pdfFrame.style.display = 'none';
        if (pdfLoadingIndicator) pdfLoadingIndicator.style.display = 'none';
    } else {
        // Desktop: use iframe as before
        console.log('Desktop device - using iframe PDF viewer');
        
        // Reset display states
        if (pdfFrame) pdfFrame.style.display = 'block';
        if (pdfFallback) pdfFallback.style.display = 'none';
        if (pdfLoadingIndicator) pdfLoadingIndicator.style.display = 'none';
        
        // Set PDF source
        if (pdfFrame) pdfFrame.src = filePath;
        
        // Set up load event handlers
        if (pdfFrame) {
            pdfFrame.onload = function() {
                console.log('PDF loaded successfully on desktop');
            };
            
            pdfFrame.onerror = function() {
                console.log('PDF failed to load on desktop');
                showPdfFallback();
            };
        }
    }
    
    if (pdfViewer) pdfViewer.style.display = 'block';
}

// Helper function to show PDF fallback
function showPdfFallback() {
    console.log('Showing PDF fallback');
    const pdfFrame = document.getElementById('pdfFrame');
    const pdfFallback = document.querySelector('.pdf-fallback');
    
    if (pdfFrame) pdfFrame.style.display = 'none';
    if (pdfFallback) pdfFallback.style.display = 'block';
}

// Enhanced mobile PDF detection
function isMobileDevice() {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) || 
           window.innerWidth <= 768;
}

// Check if browser supports PDF viewing
function supportsPdfViewing() {
    // Check if the browser supports PDF viewing in iframes
    const testFrame = document.createElement('iframe');
    testFrame.style.display = 'none';
    document.body.appendChild(testFrame);
    
    try {
        testFrame.src = 'data:application/pdf;base64,JVBERi0xLjQKJcOkw7zDtsO8DQoxIDAgb2JqDQo8PA0KL1R5cGUgL0NhdGFsb2cNCi9QYWdlcyAyIDAgUg0KPj4NCmVuZG9iag0KMiAwIG9iag0KPDwNCi9UeXBlIC9QYWdlcw0KL0NvdW50IDANCi9LaWRzIFtdDQo+Pg0KZW5kb2JqDQp4cmVmDQowIDMNCjAwMDAwMDAwMDAgNjU1MzUgZg0KMDAwMDAwMDAwMCAwMDAwMCBuDQowMDAwMDAwMDAxIDAwMDAwIG4NCnRyYWlsZXINCjw8DQovU2l6ZSAzDQovUm9vdCAxIDAgUg0KL0luZm8gMyAwIFINCj4+DQpzdGFydHhyZWYNCjANCiUlRU9GDQo=';
        
        // If the iframe loads successfully, PDF viewing is supported
        return true;
    } catch (e) {
        return false;
    } finally {
        document.body.removeChild(testFrame);
    }
}

// Show document viewer (fallback)
function showDocumentViewer(filePath, fileName) {
    const documentViewer = document.getElementById('documentViewer');
    if (documentViewer) {
        documentViewer.style.display = 'block';
    }
}

// Hide all viewers
function hideAllViewers() {
    const viewers = ['imageViewer', 'videoViewer', 'pdfViewer', 'documentViewer'];
    viewers.forEach(viewerId => {
        const viewer = document.getElementById(viewerId);
        if (viewer) {
            viewer.style.display = 'none';
        }
    });
}

// File type detection functions
function isImageFile(fileName) {
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg'];
    return imageExtensions.some(ext => fileName.toLowerCase().endsWith(ext));
}

function isVideoFile(fileName) {
    const videoExtensions = ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.webm', '.mkv'];
    return videoExtensions.some(ext => fileName.toLowerCase().endsWith(ext));
}

function isPdfFile(fileName) {
    return fileName.toLowerCase().endsWith('.pdf');
}

function getFileType(fileName) {
    if (isImageFile(fileName)) return 'image';
    if (isVideoFile(fileName)) return 'video';
    if (isPdfFile(fileName)) return 'pdf';
    return 'document';
}

function getVideoMimeType(filePath) {
    const extension = filePath.split('.').pop().toLowerCase();
    const mimeTypes = {
        'mp4': 'video/mp4',
        'avi': 'video/x-msvideo',
        'mov': 'video/quicktime',
        'wmv': 'video/x-ms-wmv',
        'flv': 'video/x-flv',
        'webm': 'video/webm',
        'mkv': 'video/x-matroska'
    };
    return mimeTypes[extension] || 'video/mp4';
}

// Navigation functions for gallery
function nextFile() {
    if (currentFileIndex < currentFiles.length - 1) {
        currentFileIndex++;
        const file = currentFiles[currentFileIndex];
        openFileModal(file.path, file.name, file.type);
    }
}

function previousFile() {
    if (currentFileIndex > 0) {
        currentFileIndex--;
        const file = currentFiles[currentFileIndex];
        openFileModal(file.path, file.name, file.type);
    }
}

// Keyboard navigation
document.addEventListener('keydown', function(event) {
    const modal = document.getElementById('fileModal');
    if (modal && modal.classList.contains('show')) {
        switch(event.key) {
            case 'ArrowLeft':
                previousFile();
                break;
            case 'ArrowRight':
                nextFile();
                break;
            case 'Escape':
                hideModalFallback();
                break;
        }
    }
});



// Initialize modal when page loads
document.addEventListener('DOMContentLoaded', function() {
    console.log('File modal DOM loaded');
    
    // Test if modal elements exist
    const modal = document.getElementById('fileModal');
    const pdfViewer = document.getElementById('pdfViewer');
    const pdfFrame = document.getElementById('pdfFrame');
    
    console.log('Modal elements on load:', {
        modal: !!modal,
        pdfViewer: !!pdfViewer,
        pdfFrame: !!pdfFrame
    });
    
    // Add close button event listeners
    const closeButtons = document.querySelectorAll('#fileModal .btn-close, #fileModal .btn-secondary');
    closeButtons.forEach(button => {
        button.addEventListener('click', function() {
            if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
                const modal = bootstrap.Modal.getInstance(document.getElementById('fileModal'));
                if (modal) {
                    modal.hide();
                }
            } else {
                hideModalFallback();
            }
        });
    });
    
    // Add backdrop click to close
    if (modal) {
        modal.addEventListener('click', function(event) {
            if (event.target === modal) {
                if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
                    const bootstrapModal = bootstrap.Modal.getInstance(modal);
                    if (bootstrapModal) {
                        bootstrapModal.hide();
                    }
                } else {
                    hideModalFallback();
                }
            }
        });
    }
    
    // Add navigation buttons if there are multiple files
    const modalBody = document.querySelector('#fileModal .modal-body');
    if (modalBody && currentFiles.length > 1) {
        const navButtons = document.createElement('div');
        navButtons.className = 'file-navigation position-absolute';
        navButtons.style.cssText = 'top: 50%; transform: translateY(-50%); z-index: 1050;';
        navButtons.innerHTML = `
            <button class="btn btn-light btn-sm me-2" onclick="previousFile()" ${currentFileIndex === 0 ? 'disabled' : ''} aria-label="Arquivo anterior">
                <i class="bi bi-chevron-left"></i>
            </button>
            <button class="btn btn-light btn-sm" onclick="nextFile()" ${currentFileIndex === currentFiles.length - 1 ? 'disabled' : ''} aria-label="Próximo arquivo">
                <i class="bi bi-chevron-right"></i>
            </button>
        `;
        modalBody.appendChild(navButtons);
    }
    
    // Add keyboard navigation for gallery
    document.addEventListener('keydown', function(event) {
        const modal = document.getElementById('fileModal');
        if (modal && modal.classList.contains('show') && currentFiles.length > 1) {
            switch(event.key) {
                case 'ArrowLeft':
                    event.preventDefault();
                    previousFile();
                    break;
                case 'ArrowRight':
                    event.preventDefault();
                    nextFile();
                    break;
                case 'Escape':
                    event.preventDefault();
                    closeModal();
                    break;
            }
        }
    });
});

// Test function for PDF viewing
window.testPdfViewing = function() {
    console.log('Testing PDF viewing...');
    const testPdfPath = '/uploads/3db79a399469c672.pdf'; // Use an existing PDF
    openFileModal(testPdfPath, 'test.pdf', 'pdf');
};

// Function to open PDF in new tab (works better on mobile)
window.openPdfInNewTab = function(filePath) {
    console.log('Opening PDF in new tab:', filePath);
    window.open(filePath, '_blank');
};

// Export functions for global access
window.openFileModal = openFileModal;
window.showPdfViewer = showPdfViewer;
window.openPdfInNewTab = openPdfInNewTab; 
window.nextFile = nextFile;
window.previousFile = previousFile; 