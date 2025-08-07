// File Modal functionality
let currentFileIndex = 0;
let currentFiles = [];

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
        console.log('Showing image viewer');
        showImageViewer(filePath);
        modalTitle.innerHTML = '<i class="bi bi-image me-2"></i>Visualizar Imagem';
    } else if (fileType === 'video' || isVideoFile(fileName)) {
        console.log('Showing video viewer');
        showVideoViewer(filePath);
        modalTitle.innerHTML = '<i class="bi bi-camera-video me-2"></i>Visualizar Vídeo';
    } else if (isPdfFile(fileName)) {
        console.log('Showing PDF viewer');
        showPdfViewer(filePath);
        modalTitle.innerHTML = '<i class="bi bi-file-pdf me-2"></i>Visualizar PDF';
    } else {
        console.log('Showing document viewer');
        showDocumentViewer(filePath, fileName);
        modalTitle.innerHTML = '<i class="bi bi-file-earmark me-2"></i>Visualizar Documento';
    }
    
    // Try to use Bootstrap modal, fallback to manual if not available
    if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
        console.log('Using Bootstrap modal');
        const bootstrapModal = new bootstrap.Modal(modal);
        bootstrapModal.show();
    } else {
        console.log('Using fallback modal');
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

// Show image viewer
function showImageViewer(filePath) {
    const imageViewer = document.getElementById('imageViewer');
    const modalImage = document.getElementById('modalImage');
    
    modalImage.src = filePath;
    imageViewer.style.display = 'block';
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
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
}

// Check if browser supports PDF viewing
function supportsPdfViewing() {
    // Check if the browser supports PDF viewing in iframes
    const testFrame = document.createElement('iframe');
    testFrame.style.display = 'none';
    document.body.appendChild(testFrame);
    
    // Try to load a test PDF
    testFrame.src = 'data:application/pdf;base64,JVBERi0xLjQKJcOkw7zDtsO8DQoxIDAgb2JqDQo8PA0KL1R5cGUgL0NhdGFsb2cNCi9QYWdlcyAyIDAgUg0KPj4NCmVuZG9iag0KMiAwIG9iag0KPDwNCi9UeXBlIC9QYWdlcw0KL0NvdW50IDANCi9LaWRzIFtdDQo+Pg0KZW5kb2JqDQp4cmVmDQowIDMNCjAwMDAwMDAwMDAgNjU1MzUgZiANCjAwMDAwMDAwMTAgMDAwMDAgbiANCjAwMDAwMDAwNzkgMDAwMDAgbiANCnRyYWlsZXINCjw8DQovU2l6ZSAzDQovUm9vdCAxIDAgUg0KL0luZm8gMyAwIFINCj4+DQpzdGFydHhyZWYNCjExMg0KJSVFT0Y=';
    
    return new Promise((resolve) => {
        setTimeout(() => {
            const supported = testFrame.contentDocument && testFrame.contentDocument.body;
            document.body.removeChild(testFrame);
            resolve(supported);
        }, 100);
    });
}

// Show document viewer (fallback)
function showDocumentViewer(filePath, fileName) {
    const documentViewer = document.getElementById('documentViewer');
    const downloadLink = document.getElementById('downloadLink');
    
    downloadLink.href = filePath;
    documentViewer.style.display = 'block';
}

// Hide all viewers
function hideAllViewers() {
    const viewers = ['imageViewer', 'videoViewer', 'pdfViewer', 'documentViewer'];
    viewers.forEach(viewerId => {
        document.getElementById(viewerId).style.display = 'none';
    });
}

// File type detection functions
function isImageFile(fileName) {
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp'];
    return imageExtensions.some(ext => fileName.toLowerCase().endsWith(ext));
}

function isVideoFile(fileName) {
    const videoExtensions = ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.webm'];
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
    const ext = filePath.split('.').pop().toLowerCase();
    const mimeTypes = {
        'mp4': 'video/mp4',
        'avi': 'video/avi',
        'mov': 'video/quicktime',
        'wmv': 'video/x-ms-wmv',
        'flv': 'video/x-flv',
        'webm': 'video/webm'
    };
    return mimeTypes[ext] || 'video/mp4';
}

// Navigation functions for gallery
function nextFile() {
    if (currentFiles.length > 0 && currentFileIndex < currentFiles.length - 1) {
        currentFileIndex++;
        const file = currentFiles[currentFileIndex];
        openFileModal(file.path, file.name, file.type);
    }
}

function previousFile() {
    if (currentFiles.length > 0 && currentFileIndex > 0) {
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
            <button class="btn btn-light btn-sm me-2" onclick="previousFile()" ${currentFileIndex === 0 ? 'disabled' : ''}>
                <i class="bi bi-chevron-left"></i>
            </button>
            <button class="btn btn-light btn-sm" onclick="nextFile()" ${currentFileIndex === currentFiles.length - 1 ? 'disabled' : ''}>
                <i class="bi bi-chevron-right"></i>
            </button>
        `;
        modalBody.appendChild(navButtons);
    }
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