// Comments functionality
let currentCommentOffset = 10; // Start with 10 comments loaded
let currentCommentSort = 'recent';

// Initialize comments functionality
document.addEventListener('DOMContentLoaded', function() {
    initializeComments();
});

function initializeComments() {
    // Character counter for comment textarea
    const commentTextarea = document.getElementById('commentContent');
    const charCount = document.getElementById('commentCharCount');
    
    if (commentTextarea && charCount) {
        commentTextarea.addEventListener('input', function() {
            const length = this.value.length;
            charCount.textContent = length;
            
            if (length > 450) {
                charCount.style.color = '#dc3545';
            } else if (length > 400) {
                charCount.style.color = '#ffc107';
            } else {
                charCount.style.color = '#6c757d';
            }
        });
    }

    // Comment form submission
    const commentForm = document.getElementById('commentForm');
    if (commentForm) {
        commentForm.addEventListener('submit', function(e) {
            e.preventDefault();
            handleCommentSubmission();
        });
    }

    // Load more comments button
    const loadMoreBtn = document.getElementById('loadMoreComments');
    if (loadMoreBtn) {
        loadMoreBtn.addEventListener('click', function() {
            loadMoreComments();
        });
    }




}

function handleCommentSubmission() {
    const content = document.getElementById('commentContent').value.trim();
    const reportID = document.getElementById('commentReportID').value;

    if (!content) {
        showCommentError('Por favor, adicione um comentário.');
        return;
    }

    if (content.length > 500) {
        showCommentError('O comentário excede o limite de 500 caracteres.');
        return;
    }

    // Add a small delay to ensure DOM is ready
    setTimeout(() => {
        showCommentVerificationModal(reportID);
    }, 100);
}

function showCommentVerificationModal(reportId) {
    // Check if the modal element exists
    const modalElement = document.getElementById('voteVerificationModal');
    if (!modalElement) {
        console.error('Vote verification modal not found in DOM');
        showCommentError('Erro: Modal de verificação não encontrado. Recarregue a página.');
        return;
    }
    
    // Check if Bootstrap is available
    if (typeof bootstrap === 'undefined' || !bootstrap.Modal) {
        console.error('Bootstrap Modal not available');
        showCommentError('Erro: Bootstrap não carregado. Recarregue a página.');
        return;
    }
    
    // Set the report ID in the modal
    const reportIdElement = document.getElementById('voteReportId');
    if (reportIdElement) {
        reportIdElement.value = reportId;
    }
    
    // Clear previous form data
    const cpfElement = document.getElementById('voteCpf');
    const birthDateElement = document.getElementById('voteBirthDate');
    if (cpfElement) cpfElement.value = '';
    if (birthDateElement) birthDateElement.value = '';
    
    // Hide verification status
    const statusElement = document.getElementById('voteVerificationStatus');
    if (statusElement) {
        statusElement.style.display = 'none';
    }
    
    // Update modal title and content for comments
    const titleElement = document.getElementById('voteVerificationModalLabel');
    if (titleElement) {
        titleElement.innerHTML = '<i class="bi bi-shield-check me-2"></i>Verificação para Comentar';
    }
    
    const alertElement = document.querySelector('#voteVerificationModal .alert-info');
    if (alertElement) {
        alertElement.innerHTML = '<i class="bi bi-info-circle me-2"></i>Para adicionar um comentário, precisamos verificar sua identidade.';
    }
    
    const submitBtn = document.getElementById('submitVoteBtn');
    if (submitBtn) {
        submitBtn.innerHTML = '<i class="bi bi-check-circle me-1"></i>Confirmar Comentário';
        submitBtn.onclick = submitComment;
    }
    
    // Show the modal
    try {
        const modal = new bootstrap.Modal(modalElement);
        modal.show();
    } catch (error) {
        console.error('Error initializing modal:', error);
        showCommentError('Erro ao abrir modal de verificação. Recarregue a página.');
    }
}

function submitComment() {
    const reportIdElement = document.getElementById('voteReportId');
    const cpfElement = document.getElementById('voteCpf');
    const birthDateElement = document.getElementById('voteBirthDate');
    const contentElement = document.getElementById('commentContent');
    
    if (!reportIdElement || !cpfElement || !birthDateElement || !contentElement) {
        showCommentError('Erro: Elementos do formulário não encontrados.');
        return;
    }
    
    const reportId = reportIdElement.value;
    const cpf = cpfElement.value;
    const birthDate = birthDateElement.value;
    const content = contentElement.value.trim();
    
    // Validate form
    if (!cpf || !birthDate) {
        showCommentVerificationStatus('Por favor, preencha todos os campos.', 'error');
        return;
    }
    
    if (!content) {
        showCommentVerificationStatus('Por favor, adicione um comentário.', 'error');
        return;
    }
    
    // Show loading state
    showCommentVerificationStatus('Verificando CPF...', 'loading');
    document.getElementById('submitVoteBtn').disabled = true;
    
    // First verify CPF, then create comment
    fetch('/api/verify-cpf', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            cpf: cpf,
            birth_date: birthDate
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success && data.valid) {
            showCommentVerificationStatus('CPF verificado com sucesso!', 'success');
            
            // Create comment after successful verification
            setTimeout(() => {
                createCommentWithVerifiedCPF(reportId, cpf, birthDate, content);
            }, 1000);
        } else {
            showCommentVerificationStatus('CPF inválido ou data de nascimento incorreta.', 'error');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showCommentVerificationStatus('Erro ao verificar CPF. Tente novamente.', 'error');
    })
    .finally(() => {
        document.getElementById('submitVoteBtn').disabled = false;
    });
}

function createCommentWithVerifiedCPF(reportId, cpf, birthDate, content) {
    fetch('/api/comments', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            report_id: parseInt(reportId),
            cpf: cpf,
            birth_date: birthDate,
            content: content
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('voteVerificationModal'));
            modal.hide();
            
            // Clear form
            document.getElementById('commentContent').value = '';
            document.getElementById('commentCharCount').textContent = '0';
            
            // Reload comments
            currentCommentOffset = 0;
            loadComments(true);
            
            // Show success message
            showCommentSuccess('Comentário adicionado com sucesso!');
        } else {
            showCommentError(data.message || 'Erro ao adicionar comentário.');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showCommentError('Erro ao adicionar comentário. Tente novamente.');
    });
}

function createComment() {
    const content = document.getElementById('commentContent').value.trim();
    const reportID = document.getElementById('commentReportID').value;

    fetch('/api/comments', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            report_id: parseInt(reportID),
            cpf: commentCPF,
            birth_date: commentBirthDate,
            content: content
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('commentVerificationModal'));
            modal.hide();
            
            // Clear form
            document.getElementById('commentContent').value = '';
            document.getElementById('commentCharCount').textContent = '0';
            
            // Reload comments
            currentCommentOffset = 0;
            loadComments(true);
            
            // Show success message
            showCommentSuccess('Comentário adicionado com sucesso!');
        } else {
            showCommentError(data.message || 'Erro ao adicionar comentário.');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showCommentError('Erro ao adicionar comentário. Tente novamente.');
    });
}

function loadComments(replace = false) {
    const reportID = document.getElementById('commentReportID').value;
    const commentsList = document.getElementById('commentsList');
    
    // Build URL with parameters
    let url = `/api/comments?report_id=${reportID}&sort=recent`;
    
    if (!replace) {
        url += `&offset=${currentCommentOffset}`;
    }

    fetch(url)
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            if (replace) {
                commentsList.innerHTML = '';
                currentCommentOffset = 0;
            }
            
            if (data.comments && data.comments.length > 0) {
                data.comments.forEach(comment => {
                    const commentElement = createCommentElement(comment);
                    commentsList.appendChild(commentElement);
                });
                
                currentCommentOffset += data.comments.length;
            }
            
            // Update load more button
            updateLoadMoreButton(data.total);
        }
    })
    .catch(error => {
        console.error('Error loading comments:', error);
    });
}

function loadMoreComments() {
    loadComments(false);
}

function createCommentElement(comment) {
    const div = document.createElement('div');
    div.className = 'comment-item';
    div.setAttribute('data-comment-id', comment.id);
    
    div.innerHTML = `
        <div class="comment-header">
            <div class="comment-author">
                <i class="bi bi-eye-fill text-muted me-1"></i>
                <span class="author-name">OlhoUrbano${comment.hashed_cpf_display}</span>
            </div>
            <div class="comment-meta">
                <small class="text-muted">${formatDate(comment.created_at)}</small>
            </div>
        </div>
        <div class="comment-content">
            <p class="mb-2">${escapeHtml(comment.content)}</p>
        </div>
    `;
    
    return div;
}



function updateLoadMoreButton(total) {
    const loadMoreBtn = document.getElementById('loadMoreComments');
    if (loadMoreBtn) {
        if (currentCommentOffset >= total) {
            loadMoreBtn.style.display = 'none';
        } else {
            loadMoreBtn.style.display = 'block';
            loadMoreBtn.textContent = `Carregar mais comentários (${total - currentCommentOffset} restantes)`;
        }
    }
}

function showCommentVerificationStatus(message, type) {
    const statusDiv = document.getElementById('voteVerificationStatus');
    if (!statusDiv) {
        console.error('Vote verification status element not found');
        return;
    }
    
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
                <i class="bi bi-exclamation-triangle-fill text-danger me-2"></i>
                <span id="voteVerificationMessage">${message}</span>
            </div>
        `;
    }
}

function showCommentSuccess(message) {
    // Create a temporary success message
    const successDiv = document.createElement('div');
    successDiv.className = 'alert alert-success alert-dismissible fade show';
    successDiv.innerHTML = `
        <i class="bi bi-check-circle me-2"></i>
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    const commentsSection = document.querySelector('.comments-section');
    commentsSection.insertBefore(successDiv, commentsSection.firstChild);
    
    // Auto-remove after 5 seconds
    setTimeout(() => {
        successDiv.remove();
    }, 5000);
}

function showCommentError(message) {
    // Create a temporary error message
    const errorDiv = document.createElement('div');
    errorDiv.className = 'alert alert-danger alert-dismissible fade show';
    errorDiv.innerHTML = `
        <i class="bi bi-exclamation-triangle me-2"></i>
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    const commentsSection = document.querySelector('.comments-section');
    commentsSection.insertBefore(errorDiv, commentsSection.firstChild);
    
    // Auto-remove after 5 seconds
    setTimeout(() => {
        errorDiv.remove();
    }, 5000);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
