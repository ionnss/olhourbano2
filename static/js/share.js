// Share functionality for report detail page
let shareModal = null;

// Initialize share functionality
function initShare() {
    // Create share modal if it doesn't exist
    if (!shareModal) {
        createShareModal();
    }
}

// Create the share modal
function createShareModal() {
    const modal = document.createElement('div');
    modal.className = 'share-modal';
    modal.id = 'shareModal';
    
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const reportUrl = `https://olhourbano.com.br/report/${reportId}`;
    
    modal.innerHTML = `
        <div class="share-modal-content" style="
            max-width: 95vw;
            max-height: 90vh;
            width: 400px;
            margin: 20px;
            overflow-y: auto;
            border-radius: 16px;
        ">
            <button class="share-modal-close" onclick="closeShareModal()" style="
                position: absolute;
                top: 15px;
                right: 15px;
                background: #f8f9fa;
                border: none;
                border-radius: 50%;
                width: 40px;
                height: 40px;
                display: flex;
                align-items: center;
                justify-content: center;
                cursor: pointer;
                font-size: 18px;
                color: #666;
                z-index: 10;
            ">
                <i class="bi bi-x"></i>
            </button>
            
            <div class="share-modal-header" style="
                padding: 25px 20px 15px 20px;
                text-align: center;
            ">
                <h4 style="
                    margin: 0 0 8px 0;
                    font-size: 20px;
                    font-weight: 600;
                    color: #333;
                ">Compartilhar Denúncia</h4>
                <p style="
                    margin: 0;
                    font-size: 14px;
                    color: #666;
                    line-height: 1.4;
                ">Escolha como você quer compartilhar esta denúncia</p>
            </div>
            
            <!-- Link Display Box with Copy Functionality -->
            <div style="
                background: white;
                padding: 20px;
                border-radius: 12px;
                border: 1px solid #e9ecef;
                margin: 0 20px 20px 20px;
            ">
                <div style="
                    font-size: 13px;
                    color: #666;
                    margin-bottom: 10px;
                    font-weight: 500;
                ">Link da Denúncia:</div>
                <div style="
                    display: flex;
                    align-items: center;
                    gap: 10px;
                ">
                    <div style="
                        flex: 1;
                        font-size: 13px;
                        color: #333;
                        word-break: break-all;
                        font-family: 'Courier New', monospace;
                        background: #f8f9fa;
                        padding: 12px 14px;
                        border-radius: 8px;
                        border: none;
                        outline: none;
                        cursor: text;
                        user-select: all;
                        min-height: 44px;
                        display: flex;
                        align-items: center;
                    " onclick="copyShareLink()" title="Clique para copiar">${reportUrl}</div>
                    <button style="
                        background: white;
                        color: #666;
                        border: 1px solid #dee2e6;
                        padding: 12px 14px;
                        border-radius: 8px;
                        cursor: pointer;
                        font-size: 14px;
                        font-weight: 500;
                        white-space: nowrap;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        min-width: 48px;
                        height: 48px;
                        flex-shrink: 0;
                    " onclick="copyShareLink()" title="Copiar link">
                        <i class="bi bi-clipboard" style="font-size: 16px;"></i>
                    </button>
                </div>
            </div>
            
            <div class="share-preview" id="sharePreview" style="
                padding: 0 20px 20px 20px;
            ">
                <p style="
                    text-align: center;
                    color: #666;
                    font-size: 14px;
                    margin: 0;
                ">Gerando preview...</p>
            </div>
        </div>
    `;
    
    document.body.appendChild(modal);
    shareModal = modal;
}

// Show share modal
function showShareModal() {
    if (!shareModal) {
        createShareModal();
    }
    
    // Generate share image
    generateShareImage();
    
    // Add mobile-specific styles
    shareModal.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.6);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 10000;
        padding: 10px;
        box-sizing: border-box;
    `;
    
    document.body.style.overflow = 'hidden';
    
    // Add touch event handling for mobile (passive listener)
    shareModal.addEventListener('touchstart', function(e) {
        if (e.target === shareModal) {
            closeShareModal();
        }
    }, { passive: true });
}

// Close share modal
function closeShareModal() {
    if (shareModal) {
        shareModal.style.display = 'none';
        document.body.style.overflow = 'auto';
        
        // Remove touch event listener
        shareModal.removeEventListener('touchstart', function(e) {
            if (e.target === shareModal) {
                closeShareModal();
            }
        });
    }
}

async function generateShareImage() {
    const previewElement = document.getElementById('sharePreview');
    if (!previewElement) return;
    
    previewElement.innerHTML = '<p>Gerando preview...</p>';
    
    try {
        // Get report data
        const reportId = document.querySelector('.share-btn').dataset.reportId;
        const categoryName = document.querySelector('.category-name').textContent;
        const categoryIcon = document.querySelector('.category-icon').textContent;
        const description = document.querySelector('.description-text').textContent;
        const location = document.querySelector('.location-text').textContent;
        const voteCount = document.querySelector('.vote-count').textContent;
        const createdAt = document.querySelector('.report-date span').textContent;
        const reportUrl = `https://olhourbano.com.br/report/${reportId}`;
        
        // Check if html2canvas is available
        if (typeof html2canvas === 'undefined') {
            throw new Error('html2canvas library not loaded');
        }
        
        // Check if we're on mobile and html2canvas might have issues
        const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
        const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
        
        if (isMobile && isIOS) {
            // On iOS, html2canvas can be problematic, show fallback
            console.log('iOS detected, using fallback preview');
            showFallbackPreview();
            return;
        }
        
        // Generate both landscape and portrait versions
        let landscapeBlob, portraitBlob;
        
        try {
            landscapeBlob = await generateLandscapeImage(reportId, categoryName, categoryIcon, description, location, voteCount, createdAt, reportUrl);
            console.log('Landscape image generated successfully');
        } catch (error) {
            console.error('Error generating landscape image:', error);
            if (isMobile) {
                console.log('Mobile device detected, showing fallback preview');
                showFallbackPreview();
                return;
            }
            throw error;
        }
        
        try {
            portraitBlob = await generatePortraitImage(reportId, categoryName, categoryIcon, description, location, voteCount, createdAt, reportUrl);
            console.log('Portrait image generated successfully');
        } catch (error) {
            console.error('Error generating portrait image:', error);
            if (isMobile) {
                console.log('Mobile device detected, showing fallback preview');
                showFallbackPreview();
                return;
            }
            throw error;
        }
        
        // Store both versions
        window.shareImageBlob = landscapeBlob;
        window.shareImageBlobPortrait = portraitBlob;
        
        // Create URLs for both
        window.shareImageUrl = URL.createObjectURL(landscapeBlob);
        window.shareImageUrlPortrait = URL.createObjectURL(portraitBlob);
        
        console.log('Generated images:', { landscape: landscapeBlob, portrait: portraitBlob });
        
        // Show both versions in preview
        previewElement.innerHTML = `
            <div style="margin-bottom: 20px;">
                <h5 style="margin-bottom: 10px; color: #333;">Versão Paisagem (1200x630)</h5>
                <img src="${window.shareImageUrl}" alt="Preview da denúncia - Paisagem" style="max-width: 100%; height: auto; border-radius: 8px; border: 2px solid #e9ecef;">
            </div>
            <div style="margin-bottom: 20px;">
                <h5 style="margin-bottom: 10px; color: #333;">Versão Retrato (630x1200)</h5>
                <img src="${window.shareImageUrlPortrait}" alt="Preview da denúncia - Retrato" style="max-width: 100%; height: auto; border-radius: 8px; border: 2px solid #e9ecef;">
            </div>
            <div style="margin-top: 15px; display: flex; gap: 10px; justify-content: center; flex-wrap: wrap;">
                <button onclick="downloadShareImage('landscape')" style="
                    background: #28a745;
                    color: white;
                    border: none;
                    padding: 8px 16px;
                    border-radius: 6px;
                    font-size: 14px;
                    cursor: pointer;
                    display: flex;
                    align-items: center;
                    gap: 5px;
                ">
                    <i class="bi bi-download"></i>
                    Baixar Paisagem
                </button>
                <button onclick="downloadShareImage('portrait')" style="
                    background: #17a2b8;
                    color: white;
                    border: none;
                    padding: 8px 16px;
                    border-radius: 6px;
                    font-size: 14px;
                    cursor: pointer;
                    display: flex;
                    align-items: center;
                    gap: 5px;
                ">
                    <i class="bi bi-download"></i>
                    Baixar Retrato
                </button>
                <button onclick="copyImageToClipboard()" style="
                    background: #007bff;
                    color: white;
                    border: none;
                    padding: 8px 16px;
                    border-radius: 6px;
                    font-size: 14px;
                    cursor: pointer;
                    display: flex;
                    align-items: center;
                    gap: 5px;
                ">
                    <i class="bi bi-clipboard"></i>
                    Copiar Imagem
                </button>
            </div>
            <p style="margin-top: 10px; font-size: 14px; color: #666; text-align: center;">
                Baixe a imagem desejada para compartilhar nas redes sociais
            </p>
        `;
        
    } catch (error) {
        console.error('Error generating share images:', error);
        showFallbackPreview();
    }
}

async function generateLandscapeImage(reportId, categoryName, categoryIcon, description, location, voteCount, createdAt, reportUrl) {
    // Create a temporary element for the landscape share image
    const shareImageElement = document.createElement('div');
    shareImageElement.style.cssText = `
        width: 1200px;
        height: 630px;
        background: #f8f9fa;
        padding: 40px;
        box-sizing: border-box;
        font-family: 'Montserrat', sans-serif;
        color: #333;
        position: relative;
        overflow: hidden;
    `;
    
    shareImageElement.innerHTML = `
        <div style="
            position: relative;
            z-index: 1;
            height: 100%;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
        ">
            <!-- Header with Logo -->
            <div style="display: flex; justify-content: space-between; align-items: flex-start;">
                <div style="display: flex; align-items: flex-start; gap: 12px; flex: 1; min-width: 0;">
                    <div style="
                        font-size: 26px;
                        color: #333;
                        font-family: 'Segoe UI Emoji', 'Apple Color Emoji', 'Noto Color Emoji', sans-serif;
                        margin-top: 5px;
                        flex-shrink: 0;
                    ">${categoryIcon}</div>
                    <div style="min-width: 0; flex: 1;">
                        <h1 style="margin: 0; font-size: 30px; font-weight: 700; color: #333; line-height: 1.2; word-wrap: break-word; overflow-wrap: break-word;">${categoryName}</h1>
                        <p style="margin: 3px 0 0 0; font-size: 15px; color: #666;">Denúncia #${reportId}</p>
                        <p style="margin: 3px 0 0 0; font-size: 13px; color: #666;">${createdAt}</p>
                    </div>
                </div>
                <div style="
                    background: #28a745;
                    color: white;
                    padding: 10px 16px;
                    border-radius: 50px;
                    text-align: center;
                    flex-shrink: 0;
                    margin-left: 15px;
                ">
                    <div style="font-size: 20px; font-weight: 700;">${voteCount}</div>
                    <div style="font-size: 11px; opacity: 0.9;">votos</div>
                </div>
            </div>
            
            <!-- Content -->
            <div style="flex: 1; display: flex; flex-direction: column; justify-content: center; padding: 30px 0;">
                <div style="
                    background: #f8f9fa;
                    padding: 40px;
                    border-radius: 20px;
                    margin-bottom: 15px;
                ">
                    <p style="
                        font-size: 28px;
                        line-height: 1.5;
                        margin: 0 0 30px 0;
                        font-weight: 500;
                        color: #333;
                    ">${description.length > 200 ? description.substring(0, 200) + '...' : description}</p>
                    
                    <div style="
                        display: flex;
                        align-items: center;
                        gap: 15px;
                        font-size: 20px;
                        color: #666;
                        padding: 20px;
                        background: #f8f9fa;
                        border-radius: 8px;
                    ">
                        <span style="font-size: 24px;">📍</span>
                        <span>${location.length > 100 ? location.substring(0, 100) + '...' : location}</span>
                    </div>
                </div>
            </div>
            
            <!-- Footer with Logo -->
            <div style="
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding-top: 20px;
                padding-bottom: 15px;
                border-top: 2px solid #e9ecef;
            ">
                <div style="display: flex; align-items: center; gap: 15px;">
                    <img src="/static/resource/circular_eye.png" alt="Olho Urbano" style="
                        width: 50px;
                        height: 50px;
                        border-radius: 50%;
                    ">
                    <div>
                        <div style="font-size: 20px; font-weight: 600; color: #333;">Olho Urbano</div>
                        <div style="font-size: 14px; color: #666;">Cidadania Ativa</div>
                        <div style="font-size: 14px; color: #999;">olhourbano.com.br</div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Temporarily add to DOM for html2canvas
    shareImageElement.style.position = 'absolute';
    shareImageElement.style.left = '-9999px';
    document.body.appendChild(shareImageElement);
    
    // Wait for fonts to load
    await document.fonts.ready;
    
    // Generate image using html2canvas with mobile-friendly settings
    const canvas = await html2canvas(shareImageElement, {
        width: 1200,
        height: 630,
        scale: 1,
        useCORS: true,
        allowTaint: true,
        backgroundColor: null,
        logging: false,
        removeContainer: true,
        foreignObjectRendering: false,
        imageTimeout: 15000, // 15 second timeout for images
        ignoreElements: (element) => {
            // Ignore elements that might cause issues on mobile
            return element.tagName === 'IFRAME' || 
                   element.classList.contains('mobile-ignore');
        }
    });
    
    // Remove temporary element
    document.body.removeChild(shareImageElement);
    
    return new Promise((resolve) => {
        canvas.toBlob((blob) => {
            resolve(blob);
        }, 'image/png', 0.9);
    });
}

async function generatePortraitImage(reportId, categoryName, categoryIcon, description, location, voteCount, createdAt, reportUrl) {
    // Create a temporary element for the portrait share image
    const shareImageElement = document.createElement('div');
    shareImageElement.style.cssText = `
        width: 630px;
        height: 1200px;
        background: #f8f9fa;
        padding: 40px;
        box-sizing: border-box;
        font-family: 'Montserrat', sans-serif;
        color: #333;
        position: relative;
        overflow: hidden;
    `;
    
    shareImageElement.innerHTML = `
        <div style="
            position: relative;
            z-index: 1;
            height: 100%;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
        ">
            <!-- Header with Logo -->
            <div style="display: flex; flex-direction: column; align-items: center; text-align: center; gap: 15px;">
                <div style="display: flex; align-items: center; gap: 15px;">
                    <div style="
                        font-size: 40px;
                        color: #333;
                        font-family: 'Segoe UI Emoji', 'Apple Color Emoji', 'Noto Color Emoji', sans-serif;
                    ">${categoryIcon}</div>
                    <div>
                        <h1 style="margin: 0; font-size: 42px; font-weight: 700; color: #333;">${categoryName}</h1>
                    </div>
                </div>
                <div style="
                    display: flex;
                    align-items: center;
                    gap: 30px;
                    justify-content: center;
                ">
                    <div style="text-align: left;">
                        <p style="margin: 0; font-size: 20px; color: #666;">Denúncia #${reportId}</p>
                        <p style="margin: 3px 0 0 0; font-size: 18px; color: #666;">${createdAt}</p>
                    </div>
                    <div style="
                        background: #28a745;
                        color: white;
                        padding: 20px 30px;
                        border-radius: 50px;
                        text-align: center;
                    ">
                        <div style="font-size: 28px; font-weight: 700;">${voteCount}</div>
                        <div style="font-size: 16px; opacity: 0.9;">votos</div>
                    </div>
                </div>
            </div>
            
            <!-- Content -->
            <div style="flex: 1; display: flex; flex-direction: column; justify-content: center; padding: 40px 0;">
                <div style="
                    background: #f8f9fa;
                    padding: 50px;
                    border-radius: 20px;
                    margin-bottom: 20px;
                ">
                    <p style="
                        font-size: 32px;
                        line-height: 1.5;
                        margin: 0 0 40px 0;
                        font-weight: 500;
                        color: #333;
                        text-align: center;
                    ">${description.length > 150 ? description.substring(0, 150) + '...' : description}</p>
                    
                    <div style="
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        gap: 15px;
                        font-size: 24px;
                        color: #666;
                        padding: 25px;
                        background: #f8f9fa;
                        border-radius: 8px;
                    ">
                        <span style="font-size: 28px;">📍</span>
                        <span style="text-align: center;">${location.length > 80 ? location.substring(0, 80) + '...' : location}</span>
                    </div>
                </div>
            </div>
            
            <!-- Footer with Logo -->
            <div style="
                display: flex;
                flex-direction: column;
                align-items: center;
                gap: 15px;
                padding-top: 20px;
                border-top: 2px solid #e9ecef;
                text-align: center;
            ">
                <div style="display: flex; align-items: center; gap: 15px;">
                    <img src="/static/resource/circular_eye.png" alt="Olho Urbano" style="
                        width: 60px;
                        height: 60px;
                        border-radius: 50%;
                    ">
                    <div>
                        <div style="font-size: 24px; font-weight: 600; color: #333;">Olho Urbano</div>
                        <div style="font-size: 16px; color: #666;">Cidadania Ativa</div>
                        <div style="font-size: 16px; color: #999;">olhourbano.com.br</div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    // Temporarily add to DOM for html2canvas
    shareImageElement.style.position = 'absolute';
    shareImageElement.style.left = '-9999px';
    document.body.appendChild(shareImageElement);
    
    // Wait for fonts to load
    await document.fonts.ready;
    
    // Generate image using html2canvas with mobile-friendly settings
    const canvas = await html2canvas(shareImageElement, {
        width: 630,
        height: 1200,
        scale: 1,
        useCORS: true,
        allowTaint: true,
        backgroundColor: null,
        logging: false,
        removeContainer: true,
        foreignObjectRendering: false,
        imageTimeout: 15000, // 15 second timeout for images
        ignoreElements: (element) => {
            // Ignore elements that might cause issues on mobile
            return element.tagName === 'IFRAME' || 
                   element.classList.contains('mobile-ignore');
        }
    });
    
    // Remove temporary element
    document.body.removeChild(shareImageElement);
    
    return new Promise((resolve) => {
        canvas.toBlob((blob) => {
            resolve(blob);
        }, 'image/png', 0.9);
    });
}


// Share functions
function shareToWhatsApp() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    const text = `Denúncia #${reportId} no Olho Urbano - Cidadania Ativa! 👁️`;
    
    const whatsappUrl = `https://wa.me/?text=${encodeURIComponent(text + '\n\n' + url)}`;
    window.open(whatsappUrl, '_blank');
}

function shareToFacebook() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    
    const facebookUrl = `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`;
    window.open(facebookUrl, '_blank');
}

function shareToTwitter() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    const text = `Denúncia #${reportId} no Olho Urbano - Cidadania Ativa! 👁️`;
    
    const twitterUrl = `https://twitter.com/intent/tweet?text=${encodeURIComponent(text)}&url=${encodeURIComponent(url)}`;
    window.open(twitterUrl, '_blank');
}

function shareToTelegram() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    const text = `Denúncia #${reportId} no Olho Urbano - Cidadania Ativa! 👁️`;
    
    const telegramUrl = `https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(text)}`;
    window.open(telegramUrl, '_blank');
}

function shareToInstagram() {
    // Create Instagram modal
    const instagramModal = document.createElement('div');
    instagramModal.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0,0,0,0.8);
        z-index: 10001;
        display: flex;
        align-items: center;
        justify-content: center;
    `;
    
    instagramModal.innerHTML = `
        <div style="
            background: white;
            padding: 2rem;
            border-radius: 16px;
            max-width: 500px;
            text-align: center;
        ">
            <h4 style="margin-bottom: 1rem; color: #333;">Compartilhar no Instagram</h4>
            
            <div style="margin-bottom: 1.5rem;">
                <p style="margin-bottom: 1rem; color: #666;">Escolha como compartilhar:</p>
                
                <div style="display: flex; flex-direction: column; gap: 1rem;">
                    <button onclick="shareToInstagramStory()" style="
                        background: linear-gradient(45deg, #f09433, #e6683c, #dc2743, #cc2366, #bc1888);
                        color: white;
                        border: none;
                        padding: 1rem;
                        border-radius: 12px;
                        font-weight: 600;
                        cursor: pointer;
                    ">
                        📱 Instagram Stories
                    </button>
                    
                    <button onclick="shareToInstagramPost()" style="
                        background: #e4405f;
                        color: white;
                        border: none;
                        padding: 1rem;
                        border-radius: 12px;
                        font-weight: 600;
                        cursor: pointer;
                    ">
                        📸 Instagram Post
                    </button>
                </div>
            </div>
            
            <button onclick="closeInstagramModal()" style="
                background: #6c757d;
                color: white;
                border: none;
                padding: 0.5rem 1rem;
                border-radius: 6px;
                cursor: pointer;
            ">
                Fechar
            </button>
        </div>
    `;
    
    document.body.appendChild(instagramModal);
    
    // Store modal reference
    window.instagramModal = instagramModal;
}

function shareToLinkedIn() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    const text = `Denúncia #${reportId} no Olho Urbano - Cidadania Ativa! 👁️`;
    
    const linkedinUrl = `https://www.linkedin.com/sharing/share-offsite/?url=${encodeURIComponent(url)}`;
    window.open(linkedinUrl, '_blank');
}

// Instagram-specific functions
function shareToInstagramStory() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    const text = `🚨 Denúncia #${reportId}\n👁️ Olho Urbano\n\nBaixe a imagem e adicione aos Stories!\n\n${url}`;
    
    // Copy text and close modal
    navigator.clipboard.writeText(text).then(() => {
        alert('Texto copiado! Cole no Instagram Stories. Baixe a imagem e adicione-a aos seus Stories.');
        closeInstagramModal();
    }).catch(() => {
        alert('Erro ao copiar. Baixe a imagem e adicione-a aos seus Instagram Stories.');
        closeInstagramModal();
    });
}

function shareToInstagramPost() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    const text = `🚨 Denúncia #${reportId} no Olho Urbano!\n\n👁️ Cidadania Ativa em ação\n\nBaixe a imagem e adicione ao seu post!\n\n#OlhoUrbano #CidadaniaAtiva #Denúncia\n\n${url}`;
    
    // Copy text and close modal
    navigator.clipboard.writeText(text).then(() => {
        alert('Texto copiado! Cole no Instagram. Baixe a imagem e adicione-a ao seu post.');
        closeInstagramModal();
    }).catch(() => {
        alert('Erro ao copiar. Baixe a imagem e adicione-a ao seu Instagram post.');
        closeInstagramModal();
    });
}

function closeInstagramModal() {
    if (window.instagramModal) {
        document.body.removeChild(window.instagramModal);
        window.instagramModal = null;
    }
}

async function copyShareLink() {
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const url = `https://olhourbano.com.br/report/${reportId}`;
    
    try {
        // Try modern clipboard API first
        if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(url);
        } else {
            // Fallback for older browsers or non-secure contexts
            const textArea = document.createElement('textarea');
            textArea.value = url;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            textArea.style.top = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
        }
        
        // Show success message
        const copyButton = document.querySelector('.share-modal button[onclick="copyShareLink()"]');
        if (copyButton) {
            const originalHTML = copyButton.innerHTML;
            copyButton.innerHTML = '<i class="bi bi-check"></i>';
            copyButton.style.background = '#28a745';
            
            setTimeout(() => {
                copyButton.innerHTML = originalHTML;
                copyButton.style.background = '#007bff';
            }, 2000);
        }
        
    } catch (error) {
        console.error('Error copying to clipboard:', error);
        
        // Show user-friendly error message
        const copyButton = document.querySelector('.share-modal button[onclick="copyShareLink()"]');
        if (copyButton) {
            const originalHTML = copyButton.innerHTML;
            copyButton.innerHTML = '<i class="bi bi-x"></i>';
            copyButton.style.background = '#dc3545';
            
            setTimeout(() => {
                copyButton.innerHTML = originalHTML;
                copyButton.style.background = '#007bff';
            }, 2000);
        }
    }
}

// Main share function called from the button
function shareReport() {
    showShareModal();
}

// Close modal when clicking outside
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('share-modal')) {
        closeShareModal();
    }
});

// Close modal with Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && shareModal && shareModal.style.display === 'block') {
        closeShareModal();
    }
});

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', initShare);

// Download share image
function downloadShareImage(version = 'landscape') {
    if (version === 'landscape') {
        if (window.shareImageBlob) {
            const url = window.shareImageUrl;
            const a = document.createElement('a');
            a.href = url;
            a.download = `denuncia-${document.querySelector('.share-btn').dataset.reportId}-paisagem.png`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        }
    } else { // portrait
        if (window.shareImageBlobPortrait) {
            const url = window.shareImageUrlPortrait;
            const a = document.createElement('a');
            a.href = url;
            a.download = `denuncia-${document.querySelector('.share-btn').dataset.reportId}-retrato.png`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        }
    }
}

// Copy image to clipboard
async function copyImageToClipboard() {
    if (window.shareImageBlob) {
        try {
            const clipboardItem = new ClipboardItem({
                'image/png': window.shareImageBlob
            });
            await navigator.clipboard.write([clipboardItem]);
            alert('Imagem copiada para a área de transferência!');
        } catch (error) {
            console.error('Error copying image:', error);
            alert('Erro ao copiar imagem. Tente baixar a imagem primeiro.');
        }
    }
}

// Fallback preview function for mobile devices
function showFallbackPreview() {
    const previewElement = document.getElementById('sharePreview');
    if (!previewElement) return;
    
    const reportId = document.querySelector('.share-btn').dataset.reportId;
    const categoryName = document.querySelector('.category-name').textContent;
    const categoryIcon = document.querySelector('.category-icon').textContent;
    const description = document.querySelector('.description-text').textContent;
    const location = document.querySelector('.location-text').textContent;
    const voteCount = document.querySelector('.vote-count').textContent;
    const createdAt = document.querySelector('.report-date span').textContent;
    const reportUrl = `https://olhourbano.com.br/report/${reportId}`;
    
    const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    
    const fallbackPreview = `
        <div style="
            background: #f8f9fa;
            color: #333;
            padding: 2rem;
            border-radius: 20px;
            text-align: left;
            font-family: 'Montserrat', sans-serif;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            border: 1px solid #e9ecef;
        ">
            ${isMobile ? `
            <div style="
                background: #e3f2fd;
                border: 1px solid #2196f3;
                color: #1565c0;
                padding: 1rem;
                border-radius: 12px;
                margin-bottom: 1.5rem;
                font-size: 0.9rem;
            ">
                <i class="bi bi-info-circle me-2"></i>
                <strong>Preview Mobile:</strong> A visualização da imagem foi simplificada para melhor compatibilidade com dispositivos móveis.
            </div>
            ` : ''}
            <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1.5rem;">
                <div style="display: flex; align-items: flex-start; gap: 1rem; flex: 1;">
                    <div style="
                        font-size: 1.5rem;
                        color: #333;
                        font-family: 'Segoe UI Emoji', 'Apple Color Emoji', 'Noto Color Emoji', sans-serif;
                        margin-top: 0.25rem;
                    ">${categoryIcon}</div>
                    <div style="min-width: 0;">
                        <h3 style="margin: 0; font-size: 1.5rem; font-weight: 700; color: #333; line-height: 1.2;">${categoryName}</h3>
                        <p style="margin: 0.25rem 0 0 0; font-size: 1rem; color: #666;">Denúncia #${reportId}</p>
                        <p style="margin: 0.25rem 0 0 0; font-size: 0.875rem; color: #666;">${createdAt}</p>
                    </div>
                </div>
                <div style="
                    background: #28a745;
                    color: white;
                    padding: 0.75rem 1.25rem;
                    border-radius: 25px;
                    text-align: center;
                    flex-shrink: 0;
                    margin-left: 1rem;
                ">
                    <div style="font-size: 1.25rem; font-weight: 700;">${voteCount}</div>
                    <div style="font-size: 0.875rem; opacity: 0.9;">votos</div>
                </div>
            </div>
            
            <div style="
                background: #f8f9fa;
                padding: 1.5rem;
                border-radius: 20px;
                margin-bottom: 1rem;
            ">
                <p style="
                    font-size: 1.125rem;
                    line-height: 1.5;
                    margin: 0 0 1rem 0;
                    font-weight: 500;
                    color: #333;
                ">${description.length > 150 ? description.substring(0, 150) + '...' : description}</p>
                
                <div style="
                    display: flex;
                    align-items: center;
                    gap: 0.75rem;
                    font-size: 1rem;
                    color: #666;
                    padding: 1rem;
                    background: #f8f9fa;
                    border-radius: 8px;
                ">
                    <span style="font-size: 1.25rem;">📍</span>
                    <span>${location.length > 80 ? location.substring(0, 80) + '...' : location}</span>
                </div>
            </div>
            
            <div style="
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding-top: 1rem;
                border-top: 2px solid #e9ecef;
            ">
                <div style="display: flex; align-items: center; gap: 0.75rem;">
                    <img src="/static/resource/circular_eye.png" alt="Olho Urbano" style="
                        width: 40px;
                        height: 40px;
                        border-radius: 50%;
                    ">
                    <div>
                        <div style="font-size: 1rem; font-weight: 600; color: #333;">Olho Urbano</div>
                        <div style="font-size: 0.875rem; color: #666;">Cidadania Ativa</div>
                        <div style="font-size: 0.75rem; color: #999;">olhourbano.com.br</div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    previewElement.innerHTML = fallbackPreview;
}

// Export functions for global access
window.shareReport = shareReport;
window.shareToWhatsApp = shareToWhatsApp;
window.shareToFacebook = shareToFacebook;
window.shareToTwitter = shareToTwitter;
window.shareToTelegram = shareToTelegram;
window.shareToInstagram = shareToInstagram;
window.shareToInstagramStory = shareToInstagramStory;
window.shareToInstagramPost = shareToInstagramPost;
window.closeInstagramModal = closeInstagramModal;
window.shareToLinkedIn = shareToLinkedIn;
window.copyShareLink = copyShareLink;
window.closeShareModal = closeShareModal;
window.downloadShareImage = downloadShareImage;
window.copyImageToClipboard = copyImageToClipboard;
window.generateShareImage = generateShareImage;

