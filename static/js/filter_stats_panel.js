// Filter and Statistics Panel functionality
// This file handles the lateral panel, mobile panel, and statistics modal

// Lateral panel functionality
function toggleLateralPanel() {
    console.log('toggleLateralPanel called');
    const lateralPanel = document.getElementById('lateral-panel');
    const cornerButton = document.getElementById('feed-corner-button');
    
    console.log('lateralPanel:', lateralPanel);
    console.log('cornerButton:', cornerButton);
    console.log('lateralPanel.style.right:', lateralPanel?.style.right);
    
    if (lateralPanel.style.right === '0px') {
        // Close panel
        console.log('Closing panel');
        lateralPanel.style.right = '-400px';
        cornerButton.style.right = '0';
    } else {
        // Open panel
        console.log('Opening panel');
        lateralPanel.style.right = '0px';
        cornerButton.style.right = '400px';
        
        // Prevent immediate closing by adding a small delay
        setTimeout(() => {
            console.log('Panel should now be open');
        }, 100);
    }
}

// Mobile panel functionality
function toggleMobilePanel() {
    const mobilePanel = document.getElementById('mobile-bottom-panel');
    const mobileOverlay = document.getElementById('mobile-overlay');
    const cornerButton = document.getElementById('feed-corner-button');
    
    if (mobilePanel.style.bottom === '0px') {
        // Close panel
        mobilePanel.style.bottom = '-100vh';
        mobileOverlay.style.display = 'none';
        cornerButton.style.display = 'block';
    } else {
        // Open panel
        mobilePanel.style.bottom = '0px';
        mobileOverlay.style.display = 'block';
        cornerButton.style.display = 'none';
    }
}

// Unified toggle function that detects screen size
function togglePanel() {
    console.log('togglePanel called, window width:', window.innerWidth);
    if (window.innerWidth <= 768) {
        console.log('Calling toggleMobilePanel');
        toggleMobilePanel();
    } else {
        console.log('Calling toggleLateralPanel');
        toggleLateralPanel();
    }
}

// Clear lateral panel filters
function clearLateralFilters() {
    document.getElementById('lateral-category').value = '';
    document.getElementById('lateral-status').value = '';
    document.getElementById('lateral-city').value = '';
    document.getElementById('lateral-sort').value = 'recent';
}

// Clear mobile panel filters
function clearMobileFilters() {
    document.getElementById('mobile-category').value = '';
    document.getElementById('mobile-status').value = '';
    document.getElementById('mobile-city').value = '';
    document.getElementById('mobile-sort').value = 'recent';
}

// Close panels when clicking outside
document.addEventListener('click', function(event) {
    const lateralPanel = document.getElementById('lateral-panel');
    const mobilePanel = document.getElementById('mobile-bottom-panel');
    const cornerButton = document.getElementById('feed-corner-button');
    const mobileOverlay = document.getElementById('mobile-overlay');
    
    console.log('Click event - target:', event.target);
    console.log('Click event - lateralPanel contains target:', lateralPanel?.contains(event.target));
    console.log('Click event - cornerButton contains target:', cornerButton?.contains(event.target));
    
    // Close lateral panel if clicking outside
    if (window.innerWidth > 768 && lateralPanel && !lateralPanel.contains(event.target) && !cornerButton.contains(event.target)) {
        console.log('Closing panel due to click outside');
        if (lateralPanel.style.right === '0px') {
            toggleLateralPanel();
        }
    }
    
    // Close mobile panel if clicking overlay
    if (window.innerWidth <= 768 && mobileOverlay && event.target === mobileOverlay) {
        toggleMobilePanel();
    }
});

// Handle window resize
window.addEventListener('resize', function() {
    const lateralPanel = document.getElementById('lateral-panel');
    const mobilePanel = document.getElementById('mobile-bottom-panel');
    const cornerButton = document.getElementById('feed-corner-button');
    const mobileOverlay = document.getElementById('mobile-overlay');
    
    if (window.innerWidth > 768) {
        // Desktop: hide mobile panel, show corner button
        mobilePanel.style.display = 'none';
        mobileOverlay.style.display = 'none';
        cornerButton.style.display = 'block';
        lateralPanel.style.right = '-400px';
        cornerButton.style.right = '0';
    } else {
        // Mobile: hide lateral panel, show corner button
        lateralPanel.style.right = '-400px';
        cornerButton.style.right = '0';
        mobilePanel.style.display = 'block';
        mobilePanel.style.bottom = '-100vh';
    }
});

// Form submission handlers
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOMContentLoaded - filter_stats_panel.js loaded');
    
    // Check if elements exist
    const lateralPanel = document.getElementById('lateral-panel');
    const mobilePanel = document.getElementById('mobile-bottom-panel');
    const cornerButton = document.getElementById('feed-corner-button');
    
    console.log('lateralPanel found:', !!lateralPanel);
    console.log('mobilePanel found:', !!mobilePanel);
    console.log('cornerButton found:', !!cornerButton);
    
    // Set progress bar width from data attribute
    const progressBars = document.querySelectorAll('.progress-bar[data-width]');
    progressBars.forEach(bar => {
        const width = bar.getAttribute('data-width');
        bar.style.width = width + '%';
    });
    
    // Lateral panel form
    const lateralForm = document.getElementById('lateral-filter-form');
    if (lateralForm) {
        lateralForm.addEventListener('submit', function(e) {
            // Form will submit normally
        });
    }
    
    // Mobile panel form
    const mobileForm = document.getElementById('mobile-filter-form');
    if (mobileForm) {
        mobileForm.addEventListener('submit', function(e) {
            // Form will submit normally
        });
    }
    
    // Initialize panel visibility based on screen size
    if (window.innerWidth <= 768) {
        // Mobile: show mobile panel, hide lateral panel
        console.log('Initializing for mobile');
        if (mobilePanel) mobilePanel.style.display = 'block';
        if (lateralPanel) lateralPanel.style.right = '-400px';
    } else {
        // Desktop: hide mobile panel, show corner button
        console.log('Initializing for desktop');
        if (mobilePanel) mobilePanel.style.display = 'none';
        if (lateralPanel) lateralPanel.style.right = '-400px';
    }
});

// Statistics modal functionality
function openStatsModal() {
    console.log('openStatsModal called');
    const modalElement = document.getElementById('statsModal');
    console.log('Modal element:', modalElement);
    
    if (modalElement) {
        const statsModal = new bootstrap.Modal(modalElement);
        console.log('Bootstrap modal instance:', statsModal);
        statsModal.show();
    } else {
        console.error('Stats modal element not found');
    }
}

// Export functions for global access
window.toggleLateralPanel = toggleLateralPanel;
window.toggleMobilePanel = toggleMobilePanel;
window.togglePanel = togglePanel;
window.clearLateralFilters = clearLateralFilters;
window.clearMobileFilters = clearMobileFilters;
window.openStatsModal = openStatsModal; 