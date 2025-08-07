// Google Maps integration
let map;
let marker;
let autocomplete;
let geocoder;
let AdvancedMarkerElement;
let PinElement;

// Global validation state
let cpfVerificationStatus = {
    verified: false,
    error: false,
    message: '',
    lastVerifiedCPF: '',
    lastVerifiedBirthDate: ''
};

// Initialize Google Maps
async function initMap() {
    const defaultLocation = { lat: -23.5505, lng: -46.6333 }; // São Paulo
    
    // Import the marker library
    const { AdvancedMarkerElement: AME, PinElement: PE } = await google.maps.importLibrary("marker");
    AdvancedMarkerElement = AME;
    PinElement = PE;
    
    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 12,
        center: defaultLocation,
        mapId: '7574fab30cf3c8137d8b0418', // Your custom Map ID linked to olhourbano_map style
        disableDefaultUI: true,
        zoomControl: false,
        mapTypeControl: false,
        scaleControl: false,
        streetViewControl: false,
        rotateControl: false,
        fullscreenControl: false
    });
    
    // Create pin element for the marker
    const pinElement = new PinElement({
        background: '#4285F4',
        borderColor: 'white',
        glyphColor: 'white',
        scale: 1.2
    });
    
    marker = new AdvancedMarkerElement({
        position: defaultLocation,
        map: map,
        content: pinElement.element,
        gmpDraggable: true
    });
    
    geocoder = new google.maps.Geocoder();
    
    // Address autocomplete
    autocomplete = new google.maps.places.Autocomplete(
        document.getElementById('location'),
        { componentRestrictions: { country: 'br' } }
    );
    
    // Event listeners
    marker.addListener('dragend', updateLocationFromMarker);
    map.addListener('click', updateMarkerPosition);
    autocomplete.addListener('place_changed', handleAddressSelect);
}

// Update location from marker position
function updateLocationFromMarker() {
    const position = marker.position;
    updateCoordinates(position.lat, position.lng);
    reverseGeocode(position.lat, position.lng);
}

// Update marker from map click
function updateMarkerPosition(event) {
    marker.position = event.latLng;
    updateCoordinates(event.latLng.lat(), event.latLng.lng());
    reverseGeocode(event.latLng.lat(), event.latLng.lng());
}

// Handle address autocomplete selection
function handleAddressSelect() {
    const place = autocomplete.getPlace();
    if (place.geometry) {
        map.setCenter(place.geometry.location);
        map.setZoom(17);
        marker.position = place.geometry.location;
        updateCoordinates(
            place.geometry.location.lat(),
            place.geometry.location.lng()
        );
    }
}

// Update coordinate inputs
function updateCoordinates(lat, lng) {
    document.getElementById('latitude').value = lat;
    document.getElementById('longitude').value = lng;
}

// Reverse geocode coordinates to address
function reverseGeocode(lat, lng) {
    geocoder.geocode({ location: { lat, lng } }, (results, status) => {
        if (status === 'OK' && results[0]) {
            document.getElementById('location').value = results[0].formatted_address;
        }
    });
}

// Document ready
document.addEventListener('DOMContentLoaded', function() {
    // CPF formatting and validation
    const cpfInput = document.getElementById('cpf');
    if (cpfInput) {
        cpfInput.addEventListener('input', function(e) {
            let value = e.target.value.replace(/\D/g, '');
            value = value.replace(/(\d{3})(\d{3})(\d{3})(\d{2})/, '$1.$2.$3-$4');
            e.target.value = value;
            
            // Reset CPF verification status when CPF changes
            resetCPFVerificationStatus();
        });

        // CPF + Birth date verification
        cpfInput.addEventListener('blur', verifyCPFWithBirthDate);
    }

    const birthDateInput = document.getElementById('birth_date');
    if (birthDateInput) {
        birthDateInput.addEventListener('blur', verifyCPFWithBirthDate);
        birthDateInput.addEventListener('input', resetCPFVerificationStatus);
    }

    // Email confirmation validation
    const emailConfirmation = document.getElementById('email_confirmation');
    if (emailConfirmation) {
        emailConfirmation.addEventListener('input', function() {
            const email = document.getElementById('email').value;
            const confirmation = this.value;
            
            if (confirmation && email !== confirmation) {
                this.setCustomValidity('Os emails não coincidem');
            } else {
                this.setCustomValidity('');
            }
        });
    }

    // Character counter for description
    const descriptionTextarea = document.getElementById('description');
    if (descriptionTextarea) {
        descriptionTextarea.addEventListener('input', function() {
            const charCount = document.getElementById('charCount');
            const length = this.value.length;
            charCount.textContent = length;
            
            if (length < 10) {
                charCount.style.color = '#dc3545';
            } else if (length > 900) {
                charCount.style.color = '#ffc107';
            } else {
                charCount.style.color = '#6c757d';
            }
        });
    }

    // File upload handling
    const fileInput = document.getElementById('files');
    if (fileInput) {
        fileInput.addEventListener('change', handleFilePreview);
    }

    // Map toggle
    const toggleMapBtn = document.getElementById('toggleMap');
    if (toggleMapBtn) {
        toggleMapBtn.addEventListener('click', toggleMap);
    }

    // Current location
    const currentLocationBtn = document.getElementById('getCurrentLocation');
    if (currentLocationBtn) {
        currentLocationBtn.addEventListener('click', getCurrentLocation);
    }

    // Form validation
    const reportForm = document.getElementById('reportForm');
    if (reportForm) {
        reportForm.addEventListener('submit', validateReportForm);
    }

    // Transport type selection
    const transportTypeSelect = document.getElementById('transport_type');
    if (transportTypeSelect) {
        transportTypeSelect.addEventListener('change', handleTransportTypeChange);
    }
});

// CPF + Birth date verification
async function verifyCPFWithBirthDate() {
    const cpf = document.getElementById('cpf').value;
    const birthDate = document.getElementById('birth_date').value;
    const statusDiv = document.getElementById('cpfVerificationStatus');
    
    if (!cpf || !birthDate) {
        return;
    }
    
    // Check if status div exists
    if (!statusDiv) {
        console.warn('CPF verification status div not found');
        return;
    }
    
    // Check if this CPF and birth date combination has already been verified
    if (cpfVerificationStatus.verified && 
        cpfVerificationStatus.lastVerifiedCPF === cpf && 
        cpfVerificationStatus.lastVerifiedBirthDate === birthDate) {
        // Already verified - just show the success status
        showCPFVerificationStatus('CPF verificado com sucesso!', 'success');
        return;
    }
    
    // If CPF or birth date changed, reset verification status
    if (cpfVerificationStatus.lastVerifiedCPF !== cpf || 
        cpfVerificationStatus.lastVerifiedBirthDate !== birthDate) {
        cpfVerificationStatus.verified = false;
        cpfVerificationStatus.error = false;
        cpfVerificationStatus.message = '';
    }
    
    // Clear previous status
    statusDiv.style.display = 'none';
    
    try {
        const response = await fetch('/api/verify-cpf', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                cpf: cpf,
                birth_date: birthDate
            })
        });
        
        const result = await response.json();
        
        if (result.success) {
            cpfVerificationStatus = {
                verified: true,
                error: false,
                message: 'CPF verificado com sucesso!',
                lastVerifiedCPF: cpf,
                lastVerifiedBirthDate: birthDate
            };
            showCPFVerificationStatus('CPF verificado com sucesso!', 'success');
        } else {
            cpfVerificationStatus = {
                verified: false,
                error: true,
                message: result.message || 'CPF inválido ou não encontrado',
                lastVerifiedCPF: '',
                lastVerifiedBirthDate: ''
            };
            showCPFVerificationStatus(result.message || 'CPF inválido ou não encontrado', 'error');
        }
    } catch (error) {
        console.error('Error verifying CPF:', error);
        cpfVerificationStatus = {
            verified: false,
            error: true,
            message: 'Erro ao verificar CPF. Tente novamente.',
            lastVerifiedCPF: '',
            lastVerifiedBirthDate: ''
        };
        showCPFVerificationStatus('Erro ao verificar CPF. Tente novamente.', 'error');
    }
}

function showCPFVerificationStatus(message, type) {
    const statusDiv = document.getElementById('cpfVerificationStatus');
    
    // Check if status div exists
    if (!statusDiv) {
        console.warn('CPF verification status div not found');
        return;
    }
    
    const bootstrapClass = getBootstrapClass(type);
    const icon = getIcon(type);
    
    statusDiv.innerHTML = `
        <div class="alert ${bootstrapClass} alert-sm mt-2">
            <i class="${icon} me-2"></i>
            ${message}
        </div>
    `;
    statusDiv.style.display = 'block';
}

function getBootstrapClass(type) {
    switch (type) {
        case 'success': return 'alert-success';
        case 'error': return 'alert-danger';
        case 'warning': return 'alert-warning';
        default: return 'alert-info';
    }
}

function getIcon(type) {
    switch (type) {
        case 'success': return 'bi bi-check-circle-fill';
        case 'error': return 'bi bi-x-circle-fill';
        case 'warning': return 'bi bi-exclamation-triangle-fill';
        default: return 'bi bi-info-circle-fill';
    }
}

// Get current location
function getCurrentLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
            function(position) {
                const lat = position.coords.latitude;
                const lng = position.coords.longitude;
                
                // Update map and marker
                const newPosition = { lat, lng };
                map.setCenter(newPosition);
                map.setZoom(16);
                marker.position = newPosition;
                
                // Update form fields
                updateCoordinates(lat, lng);
                reverseGeocode(lat, lng);
            },
            function(error) {
                console.error('Error getting location:', error);
                alert('Erro ao obter localização. Verifique se você permitiu o acesso à localização.');
            }
        );
    } else {
        alert('Geolocalização não é suportada pelo seu navegador.');
    }
}

// Toggle map visibility
function toggleMap() {
    const mapContainer = document.getElementById('mapContainer');
    const toggleBtn = document.getElementById('toggleMap');
    
    if (mapContainer.style.display === 'none') {
        mapContainer.style.display = 'block';
        toggleBtn.innerHTML = '<i class="bi bi-map me-1"></i>Fechar Mapa';
        // Trigger resize to ensure map renders correctly
        setTimeout(() => {
            if (map) {
                google.maps.event.trigger(map, 'resize');
            }
        }, 100);
    } else {
        mapContainer.style.display = 'none';
        toggleBtn.innerHTML = '<i class="bi bi-map me-1"></i>Abrir Mapa';
    }
}

// Handle file preview
function handleFilePreview() {
    const fileInput = document.getElementById('files');
    const filePreview = document.getElementById('filePreview');
    const fileList = document.getElementById('fileList');
    
    if (fileInput.files.length > 0) {
        filePreview.classList.remove('d-none');
        fileList.innerHTML = '';
        
        Array.from(fileInput.files).forEach((file, index) => {
            const fileItem = document.createElement('div');
            fileItem.className = 'd-flex align-items-center mb-2';
            
            // Get appropriate icon based on file type
            const icon = getFileIcon(file.type);
            
            fileItem.innerHTML = `
                <i class="${icon} me-2"></i>
                <span class="flex-grow-1">${file.name}</span>
                <small class="text-muted">(${(file.size / 1024 / 1024).toFixed(2)} MB)</small>
            `;
            fileList.appendChild(fileItem);
        });
        
        updateFileInputText();
    } else {
        filePreview.classList.add('d-none');
        updateFileInputText();
    }
}

// Get appropriate icon for file type
function getFileIcon(fileType) {
    if (fileType.startsWith('image/')) {
        return 'bi bi-image';
    } else if (fileType.startsWith('video/')) {
        return 'bi bi-camera-video';
    } else if (fileType === 'application/pdf') {
        return 'bi bi-file-pdf';
    } else if (fileType === 'text/plain') {
        return 'bi bi-file-text';
    } else if (fileType.includes('word') || fileType.includes('document')) {
        return 'bi bi-file-word';
    } else {
        return 'bi bi-file-earmark';
    }
}

// Update file input text
function updateFileInputText() {
    const fileInput = document.getElementById('files');
    const fileUploadText = document.querySelector('.file-upload-text');
    
    if (fileInput.files.length > 0) {
        const fileNames = Array.from(fileInput.files).map(file => file.name);
        if (fileNames.length === 1) {
            fileUploadText.textContent = fileNames[0];
        } else {
            fileUploadText.textContent = `${fileNames.length} arquivos selecionados`;
        }
    } else {
        fileUploadText.textContent = 'Selecionar Arquivos';
    }
}

// Form validation
function validateReportForm(e) {
    e.preventDefault();
    
    const submitBtn = document.getElementById('submitBtn');
    const originalText = submitBtn.innerHTML;
    
    // Show loading state
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<i class="bi bi-hourglass-split me-2"></i>Validando...';
    
    // Validate all required fields
    const validationErrors = [];
    
    // 1. Check CPF verification
    const cpf = document.getElementById('cpf').value;
    const birthDate = document.getElementById('birth_date').value;
    
    if (!cpf || !birthDate) {
        validationErrors.push('CPF e data de nascimento são obrigatórios');
    } else {
        // Check if verification is needed for current CPF/birth date combination
        const needsVerification = !cpfVerificationStatus.verified || 
                                cpfVerificationStatus.lastVerifiedCPF !== cpf || 
                                cpfVerificationStatus.lastVerifiedBirthDate !== birthDate;
        
        if (needsVerification) {
            if (cpfVerificationStatus.error && 
                cpfVerificationStatus.lastVerifiedCPF === cpf && 
                cpfVerificationStatus.lastVerifiedBirthDate === birthDate) {
                // Last verification attempt failed for these exact values
                validationErrors.push(`CPF não verificado: ${cpfVerificationStatus.message}`);
            } else {
                // Need to verify these new values
                validationErrors.push('CPF não foi verificado. Clique em "Verificar CPF" e aguarde a verificação.');
            }
        }
    }
    
    // 2. Check email fields
    const email = document.getElementById('email').value;
    const emailConfirmation = document.getElementById('email_confirmation').value;
    
    if (!email) {
        validationErrors.push('Email é obrigatório');
    } else if (!isValidEmail(email)) {
        validationErrors.push('Email inválido');
    }
    
    if (!emailConfirmation) {
        validationErrors.push('Confirmação de email é obrigatória');
    } else if (email !== emailConfirmation) {
        validationErrors.push('Emails não coincidem');
    }
    
    // 3. Check location
    const latitude = document.getElementById('latitude').value;
    const longitude = document.getElementById('longitude').value;
    
    if (!latitude || !longitude) {
        validationErrors.push('Localização é obrigatória. Use o mapa para selecionar a localização.');
    }
    
    // 4. Check description
    const description = document.getElementById('description').value;
    
    if (!description) {
        validationErrors.push('Descrição é obrigatória');
    } else if (description.length < 10) {
        validationErrors.push('Descrição deve ter pelo menos 10 caracteres');
    } else if (description.length > 1000) {
        validationErrors.push('Descrição deve ter no máximo 1000 caracteres');
    }
    
    // 5. Check transport fields if transport is required
    const transportType = document.getElementById('transport_type');
    if (transportType && transportType.value) {
        const transportFields = document.querySelectorAll('.transport-fields input[required], .transport-fields select[required]');
        transportFields.forEach(field => {
            if (!field.value.trim()) {
                validationErrors.push(`Campo de transporte "${field.getAttribute('placeholder') || field.name}" é obrigatório`);
            }
        });
    }
    
    // 6. Check file uploads - at least one file is required
    const fileInput = document.getElementById('files');
    if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
        validationErrors.push('Pelo menos um arquivo (foto, vídeo ou documento) é obrigatório para comprovar a denúncia');
    }
    
    // If there are validation errors, show them and stop submission
    if (validationErrors.length > 0) {
        showValidationErrors(validationErrors);
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
        return false;
    }
    
    // If validation passes, show success message and submit
    submitBtn.innerHTML = '<i class="bi bi-hourglass-split me-2"></i>Enviando...';
    
    // Submit the form
    const form = document.querySelector('form');
    if (form) {
        form.submit();
    }
    
    return true;
}

// Show validation errors
function showValidationErrors(errors) {
    const statusDiv = document.getElementById('cpfVerificationStatus');
    
    if (!statusDiv) {
        console.warn('Validation status div not found');
        return;
    }
    
    const errorList = errors.map(error => `<li>${error}</li>`).join('');
    
    statusDiv.innerHTML = `
        <div class="alert alert-danger alert-sm mt-2">
            <i class="bi bi-exclamation-triangle-fill me-2"></i>
            <strong>Erros de validação:</strong>
            <ul class="mb-0 mt-2">
                ${errorList}
            </ul>
        </div>
    `;
    statusDiv.style.display = 'block';
    
    // Scroll to the error message
    statusDiv.scrollIntoView({ behavior: 'smooth', block: 'center' });
}

// Email validation helper
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Reset CPF verification status
function resetCPFVerificationStatus() {
    cpfVerificationStatus = {
        verified: false,
        error: false,
        message: '',
        lastVerifiedCPF: '',
        lastVerifiedBirthDate: ''
    };
    
    // Hide any existing status messages
    const statusDiv = document.getElementById('cpfVerificationStatus');
    if (statusDiv) {
        statusDiv.style.display = 'none';
    }
}

// Transport type change handler
function handleTransportTypeChange() {
    const transportType = this.value;
    const transportFields = document.querySelectorAll('.transport-fields');
    
    // Hide all transport field sections
    transportFields.forEach(field => {
        field.style.display = 'none';
    });
    
    // Show the selected transport type fields
    if (transportType) {
        const selectedFields = document.getElementById(transportType + '-fields');
        if (selectedFields) {
            selectedFields.style.display = 'block';
        }
    }
}