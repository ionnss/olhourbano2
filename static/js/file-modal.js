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
    const modal = document.getElementById('fileModal');
    const modalTitle = document.getElementById('fileModalLabel');
    
    // Hide all content viewers
    hideAllViewers();
    
    // Show appropriate viewer based on file type
    if (fileType === 'image' || isImageFile(fileName)) {
        showImageViewer(filePath);
        modalTitle.innerHTML = '<i class="bi bi-image me-2"></i>Visualizar Imagem';
    } else if (fileType === 'video' || isVideoFile(fileName)) {
        showVideoViewer(filePath);
        modalTitle.innerHTML = '<i class="bi bi-camera-video me-2"></i>Visualizar Vídeo';
    } else if (isPdfFile(fileName)) {
        showPdfViewer(filePath);
        modalTitle.innerHTML = '<i class="bi bi-file-pdf me-2"></i>Visualizar PDF';
    } else {
        showDocumentViewer(filePath, fileName);
        modalTitle.innerHTML = '<i class="bi bi-file-earmark me-2"></i>Visualizar Documento';
    }
    
    // Try to use Bootstrap modal, fallback to manual if not available
    if (typeof bootstrap !== 'undefined' && bootstrap.Modal) {
        const bootstrapModal = new bootstrap.Modal(modal);
        bootstrapModal.show();
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
    const pdfViewer = document.getElementById('pdfViewer');
    const pdfFrame = document.getElementById('pdfFrame');
    
    pdfFrame.src = filePath;
    pdfViewer.style.display = 'block';
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
    const modal = document.getElementById('fileModal');
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