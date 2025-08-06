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



// Get color for category - unique colors for each category
function getCategoryColor(category) {
    const colors = {
        'infraestrutura_mobilidade': '#FF6B6B',      // Red
        'ciclismo': '#4ECDC4',                       // Teal
        'acessibilidade': '#FFE66D',                 // Yellow
        'redes_energeticas_iluminacao_publica': '#FF8E53', // Orange
        'limpeza': '#A8E6CF',                        // Light Green
        'saude_publica': '#45B7D1',                  // Blue
        'seguranca_publica': '#96CEB4',              // Green
        'meio_ambiente': '#FFEAA7',                  // Light Yellow
        'equipamentos_publicos': '#DDA0DD',          // Plum
        'drenagem': '#87CEEB',                       // Sky Blue
        'obras': '#F0E68C',                          // Khaki
        'corrupcao_gestao_publica': '#DC143C',       // Crimson
        'servicos_saude_publica': '#20B2AA',         // Light Sea Green
        'educacao_publica': '#9370DB',               // Medium Purple
        'transporte_publico': '#FF6347',             // Tomato
        'comercio_fiscalizacao': '#32CD32',          // Lime Green
        'outros': '#6C5CE7'                          // Purple
    };
    
    return colors[category] || colors['outros'];
}

// Show info window for a marker
function showInfoWindow(marker, report) {
    const content = createInfoWindowContent(report);
    
    infoWindow.setContent(content);
    infoWindow.open(map, marker);
    
    // Add event listener to "View Report" button
    setTimeout(() => {
        const viewBtn = document.querySelector('.view-report-btn');
        if (viewBtn) {
            viewBtn.addEventListener('click', function() {
                window.location.href = `/report/${report.id}`;
            });
        }
    }, 100);
}

// Create info window content
function createInfoWindowContent(report) {
    const statusClass = report.status === 'approved' ? 'status-approved' : 'status-pending';
    const statusText = report.status === 'approved' ? 'APROVADA' : 'PENDENTE';
    
    return `
        <div class="map-info-window">
            <div class="info-content">
                <div class="report-title">${report.category || 'Denúncia'}</div>
                <div class="report-category">
                    <i class="bi bi-tag-fill"></i>
                    <span>${report.category || 'Categoria não especificada'}</span>
                </div>
                <div class="report-status ${statusClass}">${statusText}</div>
                <div class="report-address">
                    <i class="bi bi-geo-alt-fill"></i>
                    ${report.address || 'Endereço não especificado'}
                </div>
                <div class="report-description">
                    ${report.description || 'Descrição não disponível'}
                </div>
                <button class="view-report-btn">Ver Denúncia</button>
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