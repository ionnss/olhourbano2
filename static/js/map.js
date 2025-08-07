// Global map variable
let map;
let markers = [];
let infoWindow;
let clusterer;
let AdvancedMarkerElement;
let PinElement;
let clusterMarkers = [];
let isClustered = false;

// Initialize the map when the page loads
async function initMap(retryCount = 0) {
    console.log('initMap called', retryCount > 0 ? `(retry ${retryCount})` : '');
    
    // Check if map container exists
    const mapContainer = document.getElementById('map');
    if (!mapContainer) {
        // Maximum retries to prevent infinite loops
        if (retryCount >= 50) { // 5 seconds maximum (50 * 100ms)
            console.error('Map container not found after 50 retries. Giving up.');
            return;
        }
        
        console.log(`Map container not found, retrying in 100ms... (attempt ${retryCount + 1}/50)`);
        // Retry after a short delay
        setTimeout(() => initMap(retryCount + 1), 100);
        return;
    }
    
    console.log('Map container found:', mapContainer);
    
    // Default center (Curitiba, Brazil)
    const defaultCenter = { lat: -25.428954, lng: -49.267137 };
    
    try {
        // Import the marker library
        const { AdvancedMarkerElement: AME, PinElement: PE } = await google.maps.importLibrary("marker");
        AdvancedMarkerElement = AME;
        PinElement = PE;
        
        // Create the map
        map = new google.maps.Map(mapContainer, {
            zoom: 6,
            center: defaultCenter,
            mapTypeId: google.maps.MapTypeId.ROADMAP,
            mapId: '7574fab30cf3c8137d8b0418', // Your custom Map ID linked to olhourbano_map style
            disableDefaultUI: true,
            zoomControl: false,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false
        });
        
        console.log('Map created successfully');
        
        // Create info window
        infoWindow = new google.maps.InfoWindow();
        
        // Add click listener to close info window when clicking on map
        map.addListener('click', function() {
            infoWindow.close();
        });
        
        // Add zoom and bounds changed listeners for dynamic clustering
        map.addListener('zoom_changed', handleMapChange);
        map.addListener('bounds_changed', handleMapChange);
        
        // Load reports data after a short delay to ensure MarkerClusterer is loaded
        setTimeout(() => {
            loadReportsOnMap();
        }, 500);
        
    } catch (error) {
        console.error('Error creating map:', error);
    }
}

// Ensure initMap is available globally
window.initMap = initMap;

// Fallback initialization if Google Maps API is already loaded
if (typeof google !== 'undefined' && google.maps) {
    console.log('Google Maps API already loaded, initializing map');
    initMap();
}

// Load reports and display them on the map
function loadReportsOnMap() {
    // Get filter parameters from URL
    const urlParams = new URLSearchParams(window.location.search);
    const category = urlParams.get('category');
    const status = urlParams.get('status');
    const city = urlParams.get('city');
    
    // Build API URL
    let apiUrl = '/api/reports/map';
    const params = [];
    if (category) params.push(`category=${encodeURIComponent(category)}`);
    if (status) params.push(`status=${encodeURIComponent(status)}`);
    if (city) params.push(`city=${encodeURIComponent(city)}`);
    if (params.length > 0) {
        apiUrl += '?' + params.join('&');
    }
    
    // Fetch reports data
    fetch(apiUrl)
        .then(response => response.json())
        .then(data => {
            if (data.success && data.reports) {
                displayReportsOnMap(data.reports);
            } else {
                console.error('Failed to load reports:', data.message);
            }
        })
        .catch(error => {
            console.error('Error loading reports:', error);
            // Show some sample data for demonstration
            showSampleData();
        });
}

// Filter functionality
function toggleMapFilters() {
    const filterPanel = document.getElementById('map-filter-panel');
    if (filterPanel.style.display === 'none' || filterPanel.style.display === '') {
        filterPanel.style.display = 'block';
        populateCurrentFilters();
    } else {
        filterPanel.style.display = 'none';
    }
}

function clearMapFilters() {
    document.getElementById('map-category').value = '';
    document.getElementById('map-status').value = '';
    document.getElementById('map-city').value = '';
    
    // Reload map with no filters
    loadReportsOnMap();
    
    // Update URL without parameters
    window.history.pushState({}, '', '/map');
    
    // Hide filter panel
    document.getElementById('map-filter-panel').style.display = 'none';
}

function populateCurrentFilters() {
    const urlParams = new URLSearchParams(window.location.search);
    
    const category = urlParams.get('category');
    const status = urlParams.get('status');
    const city = urlParams.get('city');
    
    if (category) document.getElementById('map-category').value = category;
    if (status) document.getElementById('map-status').value = status;
    if (city) document.getElementById('map-city').value = city;
}

// Handle filter form submission
document.addEventListener('DOMContentLoaded', function() {
    const filterForm = document.getElementById('map-filter-form');
    if (filterForm) {
        filterForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const category = document.getElementById('map-category').value;
            const status = document.getElementById('map-status').value;
            const city = document.getElementById('map-city').value;
            
            // Build URL with filters
            let url = '/map?';
            const params = [];
            if (category) params.push(`category=${encodeURIComponent(category)}`);
            if (status) params.push(`status=${encodeURIComponent(status)}`);
            if (city) params.push(`city=${encodeURIComponent(city)}`);
            
            if (params.length > 0) {
                url += params.join('&');
            } else {
                url = '/map';
            }
            
            // Update URL and reload map
            window.history.pushState({}, '', url);
            loadReportsOnMap();
            
            // Hide filter panel
            document.getElementById('map-filter-panel').style.display = 'none';
        });
    }
});

// Display reports as markers on the map
function displayReportsOnMap(reports) {
    // Clear existing markers and clusterer
    clearMarkers();
    
    console.log('Displaying reports:', reports.length);
    console.log('MarkerClusterer available:', typeof markerClusterer !== 'undefined');
    if (typeof markerClusterer !== 'undefined') {
        console.log('MarkerClusterer.MarkerClusterer available:', typeof markerClusterer.MarkerClusterer !== 'undefined');
    }
    
    if (!reports || reports.length === 0) {
        console.log('No reports to display');
        return;
    }
    
    // Create markers for each report
    reports.forEach(report => {
        if (report.latitude && report.longitude) {
            const marker = createMarker(report);
            markers.push(marker);
        }
    });
    
    // Display markers on map with dynamic clustering
    if (markers.length > 0) {
        console.log(`Displaying ${markers.length} markers with dynamic clustering`);
        applyDynamicClustering();
    }
    
    // Fit map to show all markers with padding
    if (markers.length > 0) {
        const bounds = new google.maps.LatLngBounds();
        markers.forEach(marker => {
            bounds.extend(marker.position);
        });
        
        // Add padding to bounds for better view
        const padding = { top: 50, right: 50, bottom: 50, left: 50 };
        map.fitBounds(bounds, padding);
    }
}

// Create a marker for a report
function createMarker(report) {
    const position = {
        lat: parseFloat(report.latitude),
        lng: parseFloat(report.longitude)
    };
    
    // Create pin element with category color
    const pinElement = new PinElement({
        background: getCategoryColor(report.category),
        borderColor: 'white',
        glyphColor: 'white',
        scale: 1.2
    });
    
    const marker = new AdvancedMarkerElement({
        position: position,
        map: map,
        title: report.description || 'Report',
        content: pinElement.element
    });
    
    // Add click listener to show info window
    marker.addListener('click', function() {
        showInfoWindow(marker, report);
    });
    
    return marker;
}



// Get color for category - unique colors for each category with semantic meaning
function getCategoryColor(category) {
    const colors = {
        // Infrastructure & Transport - Orange/Red tones for construction/danger
        'infraestrutura_mobilidade': '#FF5722',      // Deep Orange (construction)
        'obras': '#FF9800',                          // Orange (construction/warning)
        'transporte_publico': '#2196F3',             // Blue (reliability/transport)
        
        // Green/Nature categories - Green tones
        'ciclismo': '#4CAF50',                       // Green (bike lanes/eco)
        'meio_ambiente': '#2E7D32',                  // Dark Green (nature)
        'limpeza': '#8BC34A',                        // Light Green (clean/fresh)
        
        // Accessibility & Health - Blue/Purple for care/assistance
        'acessibilidade': '#9C27B0',                 // Purple (accessibility/inclusion)
        'saude_publica': '#F44336',                  // Red (health/emergency)
        'servicos_saude_publica': '#E91E63',         // Pink (health services)
        
        // Security & Safety - Red/Dark tones
        'seguranca_publica': '#795548',              // Brown (security/authority)
        'corrupcao_gestao_publica': '#424242',       // Dark Gray (corruption/serious)
        
        // Utilities - Yellow/Amber for energy/utilities
        'redes_energeticas_iluminacao_publica': '#FFC107', // Amber (electricity/light)
        'drenagem': '#00BCD4',                       // Cyan (water/drainage)
        
        // Public Services - Professional blues/teals
        'equipamentos_publicos': '#607D8B',          // Blue Gray (public infrastructure)
        'educacao_publica': '#3F51B5',               // Indigo (education/knowledge)
        'comercio_fiscalizacao': '#FF7043',          // Deep Orange (business/inspection)
        
        // Default - Neutral purple
        'outros': '#9E9E9E'                          // Gray (neutral/other)
    };
    
    return colors[category] || colors['outros'];
}

// Get category info (icon and name) for display
function getCategoryInfo(category) {
    const categoryMap = {
        'infraestrutura_mobilidade': { icon: '🚧', name: 'Infraestrutura e Mobilidade' },
        'ciclismo': { icon: '🚲', name: 'Ciclismo' },
        'acessibilidade': { icon: '♿', name: 'Acessibilidade' },
        'redes_energeticas_iluminacao_publica': { icon: '🔌', name: 'Redes Elétricas e/ou Iluminação Pública' },
        'limpeza': { icon: '♻️', name: 'Limpeza Urbana & Lixo' },
        'saude_publica': { icon: '🚑', name: 'Saúde Pública' },
        'seguranca_publica': { icon: '🚨', name: 'Segurança Pública' },
        'meio_ambiente': { icon: '🌳', name: 'Meio Ambiente' },
        'equipamentos_publicos': { icon: '🏚️', name: 'Estruturas Públicas' },
        'drenagem': { icon: '🌧️', name: 'Drenagem' },
        'obras': { icon: '🧱', name: 'Obras' },
        'corrupcao_gestao_publica': { icon: '🏛️', name: 'Corrupção e Má Gestão Pública' },
        'servicos_saude_publica': { icon: '🏥', name: 'Serviços de Saúde Pública' },
        'educacao_publica': { icon: '🎓', name: 'Educação Pública' },
        'transporte_publico': { icon: '🚌', name: 'Transporte Público' },
        'comercio_fiscalizacao': { icon: '🏪', name: 'Comércio e Fiscalização' },
        'outros': { icon: '❓', name: 'Outros' }
    };
    
    return categoryMap[category] || { icon: '❓', name: 'Categoria Desconhecida' };
}

// Create media gallery for info window
function createMediaGallery(photos, reportId) {
    if (!photos || photos.length === 0) {
        return '';
    }
    
    let mediaHtml = '<div class="report-media mb-3"><div class="media-preview">';
    
    // Show up to 3 photos
    const maxVisible = 3;
    const visiblePhotos = photos.slice(0, maxVisible);
    
    visiblePhotos.forEach((photo, index) => {
        const fileType = getFileTypeFromPath(photo);
        mediaHtml += `
            <div class="media-item" onclick="openFileModal('/${photo}', '${photo}', '${fileType}')" style="cursor: pointer;">
                <img src="/${photo}" alt="Evidência ${index + 1}" class="media-thumbnail" loading="lazy">
                <div class="media-overlay-hover">
                    <i class="bi bi-eye-fill"></i>
                </div>
            </div>
        `;
    });
    
    // Show "more" overlay if there are additional photos
    if (photos.length > maxVisible) {
        const remainingCount = photos.length - maxVisible;
        const allPhotos = photos.join(',');
        mediaHtml += `
            <div class="media-overlay" data-photos="${allPhotos}" data-report-id="${reportId}" onclick="openFileGallery(this)" style="cursor: pointer;">
                <span class="overlay-text">+${remainingCount}</span>
                <div class="media-overlay-hover">
                    <i class="bi bi-images"></i>
                </div>
            </div>
        `;
    }
    
    mediaHtml += '</div></div>';
    return mediaHtml;
}

// Helper function to determine file type from path
function getFileTypeFromPath(filePath) {
    const extension = filePath.split('.').pop().toLowerCase();
    
    if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(extension)) {
        return 'image';
    } else if (['mp4', 'webm', 'avi', 'mov'].includes(extension)) {
        return 'video';
    } else if (extension === 'pdf') {
        return 'pdf';
    } else {
        return 'document';
    }
}

// Show info window for a marker
function showInfoWindow(marker, report) {
    const content = createInfoWindowContent(report);
    
    infoWindow.setContent(content);
    infoWindow.open(map, marker);
    
    // Add event listeners after a short delay to ensure DOM is ready
    setTimeout(() => {
        // View Report button
        const viewBtn = document.querySelector('.view-report-btn');
        if (viewBtn) {
            viewBtn.addEventListener('click', function() {
                window.location.href = `/report/${report.id}`;
            });
        }
        
        // Vote button - manually attach event since it's dynamically created
        const voteBtn = document.querySelector('.vote-btn-map');
        if (voteBtn && typeof showVoteVerificationModal === 'function') {
            voteBtn.addEventListener('click', function(e) {
                e.preventDefault();
                const reportId = this.getAttribute('data-report-id');
                showVoteVerificationModal(reportId, this);
            });
        }
    }, 100);
}

// Create info window content
function createInfoWindowContent(report) {
    const statusClass = report.status === 'approved' ? 'status-approved' : 'status-pending';
    const statusText = report.status === 'approved' ? 'Resolvida' : 'Pendente';
    
    // Get category info - this would ideally come from the API response
    // For now, using a mapping similar to the backend
    const categoryInfo = getCategoryInfo(report.category);
    
    // Format the author display
    const authorDisplay = report.hashed_cpf ? `OlhoUrbano${report.hashed_cpf}` : 'OlhoUrbanoAnônimo';
    
    // Format date - already formatted by backend
    const reportDate = report.created_at || 'Data não disponível';
    
    return `
        <div class="map-info-window">
            <div class="map-info-header">
                <div class="report-category">
                    <span class="category-icon">${categoryInfo.icon}</span>
                    <span class="category-name">${categoryInfo.name}</span>
                </div>
                <div class="report-status">
                    <span class="status-badge status-${report.status}">${statusText}</span>
                </div>
            </div>
            
            <div class="map-info-body">
                <div class="report-location mb-2">
                    <i class="bi bi-geo-alt-fill text-muted me-1"></i>
                    <span class="location-text">${report.address || 'Endereço não especificado'}</span>
                </div>
                
                <div class="report-description mb-3">
                    <p class="description-text">${report.description || 'Descrição não disponível'}</p>
                </div>
                
                ${createMediaGallery(report.photos, report.id)}
            </div>
            
            <div class="map-info-footer">
                <div class="report-meta">
                    <div class="meta-item">
                        <i class="bi bi-calendar3 text-muted me-1"></i>
                        <span class="meta-text">${reportDate}</span>
                    </div>
                    <div class="meta-item">
                        <i class="bi bi-eye-fill text-muted me-1"></i>
                        <span class="meta-text">${authorDisplay}</span>
                    </div>
                </div>
                
                <div class="map-info-actions">
                    <button class="view-report-btn">
                        <i class="bi bi-eye-fill me-1"></i>
                        Ver Detalhes
                    </button>
                    <button class="vote-btn vote-btn-map" data-report-id="${report.id}">
                        <i class="bi bi-hand-thumbs-up-fill me-1"></i>
                        Votar
                        <span class="vote-shield">
                            <span class="vote-count">${report.vote_count || 0}</span>
                        </span>
                    </button>
                </div>
            </div>
        </div>
    `;
}

// Clear all markers from the map
function clearMarkers() {
    // Clear marker clusterer if it exists
    if (clusterer) {
        clusterer.clearMarkers();
        clusterer = null;
    }
    
    // Clear individual markers
    markers.forEach(marker => {
        marker.setMap(null);
    });
    markers = [];
    
    // Also clear cluster markers
    clearClusterMarkers();
    isClustered = false;
}

// Show sample data for demonstration (when API is not available)
function showSampleData() {
    const sampleReports = [
        {
            id: 1,
            category: 'infraestrutura_mobilidade',
            description: 'Buraco na rua que precisa ser consertado',
            address: 'Rua das Flores, 123 - Centro',
            status: 'pending',
            latitude: -25.428954,
            longitude: -49.267137
        },
        {
            id: 2,
            category: 'acessibilidade',
            description: 'Rampa de acesso quebrada',
            address: 'Av. Paulista, 456 - Batel',
            status: 'approved',
            latitude: -25.430000,
            longitude: -49.270000
        },
        {
            id: 3,
            category: 'redes_energeticas_iluminacao_publica',
            description: 'Poste de luz queimado',
            address: 'Rua XV de Novembro, 789 - Centro',
            status: 'pending',
            latitude: -25.426000,
            longitude: -49.264000
        },
        {
            id: 4,
            category: 'ciclismo',
            description: 'Ciclovia com buracos',
            address: 'Av. Sete de Setembro, 100 - Centro',
            status: 'pending',
            latitude: -25.425000,
            longitude: -49.265000
        },
        {
            id: 5,
            category: 'limpeza',
            description: 'Lixo acumulado na calçada',
            address: 'Rua das Palmeiras, 200 - Batel',
            status: 'approved',
            latitude: -25.432000,
            longitude: -49.268000
        }
    ];
    
    displayReportsOnMap(sampleReports);
}

// Handle window resize
window.addEventListener('resize', function() {
    if (map) {
        google.maps.event.trigger(map, 'resize');
    }
});



// Handle map zoom and bounds changes for dynamic clustering
function handleMapChange() {
    if (markers.length > 1) {
        // Debounce the clustering to avoid too many recalculations
        clearTimeout(window.clusterTimeout);
        window.clusterTimeout = setTimeout(() => {
            applyDynamicClustering();
        }, 300);
    }
}

// Dynamic clustering that responds to zoom and map movement
function applyDynamicClustering() {
    const zoom = map.getZoom();
    const bounds = map.getBounds();
    
    // Adjust cluster radius based on zoom level
    let clusterRadius;
    if (zoom >= 15) {
        clusterRadius = 0.001; // Very close markers at high zoom
    } else if (zoom >= 12) {
        clusterRadius = 0.005; // Close markers at medium-high zoom
    } else if (zoom >= 9) {
        clusterRadius = 0.01; // Medium distance at medium zoom
    } else {
        clusterRadius = 0.02; // Far distance at low zoom
    }
    
    // Clear existing cluster markers
    clearClusterMarkers();
    
    // Only cluster if we have multiple markers and zoom is low enough
    if (markers.length > 1 && zoom < 15) {
        const clusters = [];
        
        markers.forEach(marker => {
            const pos = marker.position;
            let addedToCluster = false;
            
            for (let cluster of clusters) {
                const clusterCenter = cluster.center;
                const distance = Math.sqrt(
                    Math.pow(pos.lat - clusterCenter.lat, 2) + 
                    Math.pow(pos.lng - clusterCenter.lng, 2)
                );
                
                if (distance < clusterRadius) {
                    cluster.markers.push(marker);
                    cluster.center = {
                        lat: (cluster.center.lat + pos.lat) / 2,
                        lng: (cluster.center.lng + pos.lng) / 2
                    };
                    addedToCluster = true;
                    break;
                }
            }
            
            if (!addedToCluster) {
                clusters.push({
                    center: pos,
                    markers: [marker]
                });
            }
        });
        
        // Hide individual markers and show clusters
        markers.forEach(marker => marker.setMap(null));
        
        clusters.forEach(cluster => {
            if (cluster.markers.length === 1) {
                // Single marker, show it normally
                cluster.markers[0].setMap(map);
            } else {
                // Multiple markers, create a cluster
                const clusterElement = document.createElement('div');
                clusterElement.innerHTML = `
                    <div style="
                        width: 40px; 
                        height: 40px; 
                        background: #326ffe; 
                        border: 2px solid white; 
                        border-radius: 50%; 
                        display: flex; 
                        align-items: center; 
                        justify-content: center; 
                        color: white; 
                        font-family: Arial, sans-serif; 
                        font-size: 14px; 
                        font-weight: bold;
                        box-shadow: 0 2px 4px rgba(0,0,0,0.3);
                        cursor: pointer;
                    ">
                        ${cluster.markers.length}
                    </div>
                `;
                
                const clusterMarker = new AdvancedMarkerElement({
                    position: cluster.center,
                    map: map,
                    content: clusterElement,
                    title: `${cluster.markers.length} reports in this area`
                });
                
                // Store cluster marker for later removal
                clusterMarkers.push(clusterMarker);
                
                // Add click listener to expand cluster
                clusterMarker.addListener('click', function() {
                    // Show all individual markers in this cluster
                    cluster.markers.forEach(marker => marker.setMap(map));
                    // Remove this cluster marker
                    clusterMarker.setMap(null);
                    // Remove from our tracking array
                    const index = clusterMarkers.indexOf(clusterMarker);
                    if (index > -1) {
                        clusterMarkers.splice(index, 1);
                    }
                });
            }
        });
        
        isClustered = true;
    } else {
        // Show all individual markers
        markers.forEach(marker => marker.setMap(map));
        isClustered = false;
    }
}

// Clear all cluster markers
function clearClusterMarkers() {
    clusterMarkers.forEach(marker => {
        marker.setMap(null);
    });
    clusterMarkers = [];
}

// Export functions for global access
window.loadReportsOnMap = loadReportsOnMap; 