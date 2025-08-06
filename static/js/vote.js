// Vote functionality for report voting
document.addEventListener('DOMContentLoaded', function() {
    // Add click event listeners to all vote buttons
    const voteButtons = document.querySelectorAll('.vote-btn');
    
    voteButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            
            const reportId = this.getAttribute('data-report-id');
            if (!reportId) {
                console.error('No report ID found for vote button');
                return;
            }
            
            // Call vote function
            voteForReport(reportId, this);
        });
    });
});

// Function to vote for a report
function voteForReport(reportId, buttonElement) {
    // Show loading state
    const originalText = buttonElement.innerHTML;
    buttonElement.innerHTML = '<i class="bi bi-hourglass-split me-1"></i> Votando...';
    buttonElement.disabled = true;
    
    // Make API call to vote
    fetch('/api/vote', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            report_id: reportId
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Update vote count on button
            const currentVoteCount = parseInt(data.vote_count) || 0;
            buttonElement.innerHTML = `<i class="bi bi-hand-thumbs-up-fill me-1"></i> Votar (${currentVoteCount})`;
            
            // Show success feedback
            showVoteFeedback(buttonElement, 'Voto registrado!', 'success');
        } else {
            // Show error message
            showVoteFeedback(buttonElement, data.message || 'Erro ao votar', 'error');
            
            // Restore original text
            buttonElement.innerHTML = originalText;
        }
    })
    .catch(error => {
        console.error('Error voting:', error);
        
        // Show error feedback
        showVoteFeedback(buttonElement, 'Erro de conexão', 'error');
        
        // Restore original text
        buttonElement.innerHTML = originalText;
    })
    .finally(() => {
        buttonElement.disabled = false;
    });
}

// Function to show vote feedback
function showVoteFeedback(buttonElement, message, type) {
    // Create feedback element
    const feedback = document.createElement('div');
    feedback.className = `vote-feedback vote-feedback-${type}`;
    feedback.textContent = message;
    feedback.style.cssText = `
        position: absolute;
        top: -40px;
        left: 50%;
        transform: translateX(-50%);
        background: ${type === 'success' ? '#28a745' : '#dc3545'};
        color: white;
        padding: 8px 12px;
        border-radius: 6px;
        font-size: 0.8rem;
        white-space: nowrap;
        z-index: 1000;
        animation: fadeInOut 2s ease-in-out;
    `;
    
    // Add animation styles if not already present
    if (!document.querySelector('#vote-feedback-styles')) {
        const style = document.createElement('style');
        style.id = 'vote-feedback-styles';
        style.textContent = `
            @keyframes fadeInOut {
                0% { opacity: 0; transform: translateX(-50%) translateY(10px); }
                20% { opacity: 1; transform: translateX(-50%) translateY(0); }
                80% { opacity: 1; transform: translateX(-50%) translateY(0); }
                100% { opacity: 0; transform: translateX(-50%) translateY(-10px); }
            }
        `;
        document.head.appendChild(style);
    }
    
    // Position the button relatively if not already
    if (getComputedStyle(buttonElement).position === 'static') {
        buttonElement.style.position = 'relative';
    }
    
    // Add feedback to button
    buttonElement.appendChild(feedback);
    
    // Remove feedback after animation
    setTimeout(() => {
        if (feedback.parentNode) {
            feedback.parentNode.removeChild(feedback);
        }
    }, 2000);
}

// Export functions for global access
window.voteForReport = voteForReport;
window.showVoteFeedback = showVoteFeedback; 