// Vote functionality for report voting with CPF verification
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
            
            // Show CPF verification modal instead of voting directly
            showVoteVerificationModal(reportId, this);
        });
    });

    // Add CPF formatting to vote modal
    const voteCpfInput = document.getElementById('voteCpf');
    if (voteCpfInput) {
        voteCpfInput.addEventListener('input', function(e) {
            let value = e.target.value.replace(/\D/g, ''); // Remove non-digits
            if (value.length <= 11) {
                // Format as XXX.XXX.XXX-XX
                value = value.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
                value = value.replace(/(\d{3})(\d{3})(\d{3})/, '$1.$2.$3');
                value = value.replace(/(\d{3})(\d{3})/, '$1.$2');
                e.target.value = value;
            }
        });
    }
});

// Function to show the vote verification modal
function showVoteVerificationModal(reportId, buttonElement) {
    // Store the button element for later use
    window.currentVoteButton = buttonElement;
    
    // Set the report ID in the modal
    document.getElementById('voteReportId').value = reportId;
    
    // Clear previous form data
    document.getElementById('voteCpf').value = '';
    document.getElementById('voteBirthDate').value = '';
    
    // Hide verification status
    document.getElementById('voteVerificationStatus').style.display = 'none';
    
    // Show the modal
    const modal = new bootstrap.Modal(document.getElementById('voteVerificationModal'));
    modal.show();
}

// Function to submit vote after CPF verification
function submitVote() {
    const reportId = document.getElementById('voteReportId').value;
    const cpf = document.getElementById('voteCpf').value;
    const birthDate = document.getElementById('voteBirthDate').value;
    
    // Validate form
    if (!cpf || !birthDate) {
        showVoteVerificationStatus('Por favor, preencha todos os campos.', 'error');
        return;
    }
    
    // Show loading state
    showVoteVerificationStatus('Verificando CPF...', 'loading');
    document.getElementById('submitVoteBtn').disabled = true;
    
    // Make API call to vote with CPF verification
    fetch('/api/vote', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            report_id: parseInt(reportId),
            cpf: cpf,
            birth_date: birthDate
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Update vote count on button
            const currentVoteCount = parseInt(data.vote_count) || 0;
            const buttonElement = window.currentVoteButton;
            
            buttonElement.innerHTML = `<i class="bi bi-hand-thumbs-up-fill me-1"></i> Votar <span class="vote-shield"><span class="vote-count">${currentVoteCount}</span></span>`;
            
            // Add animation to the shield
            const shield = buttonElement.querySelector('.vote-shield');
            if (shield) {
                shield.classList.add('vote-updated');
                setTimeout(() => {
                    shield.classList.remove('vote-updated');
                }, 300);
            }
            
            // Show success message and close modal
            showVoteVerificationStatus('Voto registrado com sucesso!', 'success');
            setTimeout(() => {
                const modal = bootstrap.Modal.getInstance(document.getElementById('voteVerificationModal'));
                modal.hide();
            }, 1500);
            
        } else {
            // Show error message
            showVoteVerificationStatus(data.message || 'Erro ao votar', 'error');
        }
    })
    .catch(error => {
        console.error('Error voting:', error);
        showVoteVerificationStatus('Erro de conexão', 'error');
    })
    .finally(() => {
        document.getElementById('submitVoteBtn').disabled = false;
    });
}

// Function to show verification status in modal
function showVoteVerificationStatus(message, type) {
    const statusDiv = document.getElementById('voteVerificationStatus');
    const messageSpan = document.getElementById('voteVerificationMessage');
    
    statusDiv.style.display = 'block';
    
    // Remove existing classes
    statusDiv.className = 'mb-3';
    
    // Add appropriate classes based on type
    if (type === 'loading') {
        statusDiv.innerHTML = `
            <div class="verification-status">
                <div class="spinner-border spinner-border-sm me-2" role="status">
                    <span class="visually-hidden">Verificando...</span>
                </div>
                <span id="voteVerificationMessage">${message}</span>
            </div>
        `;
    } else if (type === 'success') {
        statusDiv.innerHTML = `
            <div class="verification-status success">
                <i class="bi bi-check-circle-fill text-success me-2"></i>
                <span id="voteVerificationMessage">${message}</span>
            </div>
        `;
    } else if (type === 'error') {
        statusDiv.innerHTML = `
            <div class="verification-status error">
                <i class="bi bi-exclamation-circle-fill text-danger me-2"></i>
                <span id="voteVerificationMessage">${message}</span>
            </div>
        `;
    }
}

// Function to show vote feedback (for backward compatibility)
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
window.showVoteFeedback = showVoteFeedback;
window.showVoteVerificationModal = showVoteVerificationModal;
window.submitVote = submitVote; 