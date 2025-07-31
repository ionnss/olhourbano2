// Google Maps integration
let map;
let marker;
let autocomplete;
let geocoder;

// Initialize Google Maps
function initMap() {
    const defaultLocation = { lat: -23.5505, lng: -46.6333 }; // São Paulo
    
    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 12,
        center: defaultLocation
    });
    
    marker = new google.maps.Marker({
        position: defaultLocation,
        map: map,
        draggable: true
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
    const position = marker.getPosition();
    updateCoordinates(position.lat(), position.lng());
    reverseGeocode(position.lat(), position.lng());
}

// Update marker from map click
function updateMarkerPosition(event) {
    marker.setPosition(event.latLng);
    updateCoordinates(event.latLng.lat(), event.latLng.lng());
    reverseGeocode(event.latLng.lat(), event.latLng.lng());
}

// Handle address autocomplete selection
function handleAddressSelect() {
    const place = autocomplete.getPlace();
    if (place.geometry) {
        map.setCenter(place.geometry.location);
        map.setZoom(17);
        marker.setPosition(place.geometry.location);
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
        });

        // CPF + Birth date verification
        cpfInput.addEventListener('blur', verifyCPFWithBirthDate);
    }

    const birthDateInput = document.getElementById('birth_date');
    if (birthDateInput) {
        birthDateInput.addEventListener('blur', verifyCPFWithBirthDate);
    }

    // Email confirmation validation
    const emailConfirmation = document.getElementById('email_confirmation');
    if (emailConfirmation) {
        emailConfirmation.addEventListener('blur', function() {
            const email = document.getElementById('email').value;
            const confirmation = this.value;
            
            if (email && confirmation && email !== confirmation) {
                this.setCustomValidity('Emails não conferem');
                this.classList.add('is-invalid');
            } else {
                this.setCustomValidity('');
                this.classList.remove('is-invalid');
            }
        });
    }

    // Character counter for description
    const description = document.getElementById('description');
    if (description) {
        description.addEventListener('input', function() {
            document.getElementById('charCount').textContent = this.value.length;
        });
        
        // Initial count
        if (description.value) {
            document.getElementById('charCount').textContent = description.value.length;
        }
    }

    // File preview and custom upload button
    const filesInput = document.getElementById('files');
    if (filesInput) {
        filesInput.addEventListener('change', handleFilePreview);
    }
    
    // Initialize custom file input text
    updateFileInputText();

    // Get current location button
    const getCurrentLocationBtn = document.getElementById('getCurrentLocation');
    if (getCurrentLocationBtn) {
        getCurrentLocationBtn.addEventListener('click', getCurrentLocation);
    }

    // Toggle map button
    const toggleMapBtn = document.getElementById('toggleMap');
    if (toggleMapBtn) {
        toggleMapBtn.addEventListener('click', toggleMap);
    }

    // Form validation
    const reportForm = document.getElementById('reportForm');
    if (reportForm) {
        reportForm.addEventListener('submit', validateReportForm);
    }
});

// Verify CPF with birth date using CPFHub API
async function verifyCPFWithBirthDate() {
    const cpf = document.getElementById('cpf').value.replace(/\D/g, '');
    const birthDate = document.getElementById('birth_date').value;
    
    if (cpf.length !== 11 || !birthDate) {
        return;
    }

    try {
        showCPFVerificationStatus('Verificando CPF...', 'info');

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

        if (result.valid) {
            showCPFVerificationStatus('CPF verificado com sucesso', 'success');
            document.getElementById('cpf').classList.remove('is-invalid');
            document.getElementById('cpf').classList.add('is-valid');
        } else {
            showCPFVerificationStatus('CPF ou data de nascimento inválidos', 'error');
            document.getElementById('cpf').classList.remove('is-valid');
            document.getElementById('cpf').classList.add('is-invalid');
        }
    } catch (error) {
        console.error('Erro na verificação do CPF:', error);
        showCPFVerificationStatus('Erro na verificação. Tente novamente.', 'warning');
    }
}

// Show CPF verification status
function showCPFVerificationStatus(message, type) {
    // Remove existing status
    const existingStatus = document.querySelector('.cpf-verification-status');
    if (existingStatus) {
        existingStatus.remove();
    }

    // Create status element
    const statusDiv = document.createElement('div');
    statusDiv.className = `cpf-verification-status alert alert-${getBootstrapClass(type)} alert-sm mt-2`;
    statusDiv.innerHTML = `<i class="fas fa-${getIcon(type)} me-2"></i>${message}`;

    // Insert after CPF input
    const cpfInput = document.getElementById('cpf');
    cpfInput.parentNode.appendChild(statusDiv);

    // Auto-remove after 5 seconds for non-error messages
    if (type !== 'error') {
        setTimeout(() => {
            statusDiv.remove();
        }, 5000);
    }
}

function getBootstrapClass(type) {
    const classes = {
        'info': 'info',
        'success': 'success',
        'error': 'danger',
        'warning': 'warning'
    };
    return classes[type] || 'info';
}

function getIcon(type) {
    const icons = {
        'info': 'info-circle',
        'success': 'check-circle',
        'error': 'exclamation-triangle',
        'warning': 'exclamation-circle'
    };
    return icons[type] || 'info-circle';
}

// Get current location
function getCurrentLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
            function(position) {
                const lat = position.coords.latitude;
                const lng = position.coords.longitude;
                const location = { lat, lng };
                
                map.setCenter(location);
                map.setZoom(17);
                marker.setPosition(location);
                updateCoordinates(lat, lng);
                reverseGeocode(lat, lng);
                
                // Show map if hidden
                document.getElementById('mapContainer').style.display = 'block';
                const toggleBtn = document.getElementById('toggleMap');
                toggleBtn.innerHTML = '<i class="fas fa-eye-slash me-1"></i>Ocultar Mapa';
            },
            function(error) {
                alert('Erro ao obter localização: ' + error.message);
            }
        );
    } else {
        alert('Geolocalização não é suportada neste navegador.');
    }
}

// Toggle map visibility
function toggleMap() {
    const mapContainer = document.getElementById('mapContainer');
    const toggleBtn = document.getElementById('toggleMap');
    
    if (mapContainer.style.display === 'none') {
        mapContainer.style.display = 'block';
        toggleBtn.innerHTML = '<i class="fas fa-eye-slash me-1"></i>Ocultar Mapa';
        google.maps.event.trigger(map, 'resize');
    } else {
        mapContainer.style.display = 'none';
        toggleBtn.innerHTML = '<i class="fas fa-map me-1"></i>Abrir Mapa';
    }
}

// Handle file preview
function handleFilePreview() {
    const fileList = document.getElementById('fileList');
    const filePreview = document.getElementById('filePreview');
    const files = this.files;
    
    if (files.length > 0) {
        fileList.innerHTML = '';
        Array.from(files).forEach((file, index) => {
            const fileItem = document.createElement('div');
            fileItem.className = 'mb-2 p-2 border rounded bg-white d-flex justify-content-between align-items-center';
            fileItem.innerHTML = `
                <div>
                    <i class="fas fa-file me-2"></i>
                    <strong>${file.name}</strong>
                    <span class="text-muted ms-2">(${(file.size / 1024 / 1024).toFixed(2)} MB)</span>
                </div>
                <span class="badge bg-primary">${file.type}</span>
            `;
            fileList.appendChild(fileItem);
        });
        filePreview.classList.remove('d-none');
    } else {
        filePreview.classList.add('d-none');
    }
}

// Update file input text
function updateFileInputText() {
    const fileInput = document.getElementById('files');
    const fileText = document.querySelector('.file-upload-text');
    
    if (fileInput && fileText) {
        fileInput.addEventListener('change', function() {
            const fileCount = this.files.length;
            if (fileCount > 0) {
                if (fileCount === 1) {
                    fileText.textContent = `1 arquivo selecionado`;
                } else {
                    fileText.textContent = `${fileCount} arquivos selecionados`;
                }
            } else {
                fileText.textContent = 'Selecionar Arquivos';
            }
        });
    }
}

// Validate report form
function validateReportForm(e) {
    const lat = document.getElementById('latitude').value;
    const lng = document.getElementById('longitude').value;
    
    if (!lat || !lng || lat == '0' || lng == '0') {
        e.preventDefault();
        alert('Por favor, defina a localização usando o mapa ou GPS.');
        return false;
    }

    // Check CPF verification status
    const cpfInput = document.getElementById('cpf');
    if (cpfInput.classList.contains('is-invalid')) {
        e.preventDefault();
        alert('Por favor, verifique se o CPF e data de nascimento estão corretos.');
        return false;
    }
}