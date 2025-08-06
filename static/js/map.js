// Global map variable
let map;
let markers = [];
let infoWindow;

// Initialize the map when the page loads
function initMap() {
    console.log('initMap called');
    
    // Check if map container exists
    const mapContainer = document.getElementById('map');
    if (!mapContainer) {
        console.error('Map container not found');
        return;
    }
    
    console.log('Map container found:', mapContainer);
    
    // Default center (Curitiba, Brazil)
    const defaultCenter = { lat: -25.428954, lng: -49.267137 };
    
    try {
        // Create the map
        map = new google.maps.Map(mapContainer, {
            zoom: 6,
            center: defaultCenter,
            mapTypeId: google.maps.MapTypeId.ROADMAP,
            disableDefaultUI: true,
            zoomControl: false,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            styles: [
                {
                    featureType: 'poi',
                    elementType: 'labels',
                    stylers: [{ visibility: 'off' }]
                }
            ]
        });
        
        console.log('Map created successfully');
        
        // Create info window
        infoWindow = new google.maps.InfoWindow();
        
        // Load reports data
        loadReportsOnMap();
        
        // Add click listener to close info window when clicking on map
        map.addListener('click', function() {
            infoWindow.close();
        });
        
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
    // Clear existing markers
    clearMarkers();
    
    if (!reports || reports.length === 0) {
        console.log('No reports to display');
        return;
    }
    
    reports.forEach(report => {
        if (report.latitude && report.longitude) {
            const marker = createMarker(report);
            markers.push(marker);
        }
    });
    
    // Fit map to show all markers with padding
    if (markers.length > 0) {
        const bounds = new google.maps.LatLngBounds();
        markers.forEach(marker => {
            bounds.extend(marker.getPosition());
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
    
    // Choose marker icon based on category
    const icon = getMarkerIcon(report.category);
    
    const marker = new google.maps.Marker({
        position: position,
        map: map,
        icon: icon,
        title: report.description || 'Report'
    });
    
    // Add click listener to show info window
    marker.addListener('click', function() {
        showInfoWindow(marker, report);
    });
    
    return marker;
}

// Get marker icon based on category
function getMarkerIcon(category) {
    const iconBase = {
        url: 'data:image/svg+xml;charset=UTF-8,' + encodeURIComponent(`
            <svg width="32" height="32" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                <circle cx="16" cy="16" r="14" fill="${getCategoryColor(category)}" stroke="white" stroke-width="2"/>
                <circle cx="16" cy="16" r="6" fill="white"/>
            </svg>
        `),
        scaledSize: new google.maps.Size(32, 32),
        anchor: new google.maps.Point(16, 16)
    };
    
    return iconBase;
}

// Get color for category
function getCategoryColor(category) {
    const colors = {
        'infraestrutura': '#FF6B6B',
        'acessibilidade': '#4ECDC4',
        'iluminacao': '#FFE66D',
        'seguranca': '#FF8E53',
        'corrupcao': '#A8E6CF',
        'meio-ambiente': '#45B7D1',
        'obras': '#96CEB4',
        'fiscalizacao': '#FFEAA7',
        'default': '#6C5CE7'
    };
    
    return colors[category] || colors.default;
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
    markers.forEach(marker => {
        marker.setMap(null);
    });
    markers = [];
}

// Show sample data for demonstration (when API is not available)
function showSampleData() {
    const sampleReports = [
        {
            id: 1,
            category: 'infraestrutura',
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
            category: 'iluminacao',
            description: 'Poste de luz queimado',
            address: 'Rua XV de Novembro, 789 - Centro',
            status: 'pending',
            latitude: -25.426000,
            longitude: -49.264000
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

// Export functions for global access
window.loadReportsOnMap = loadReportsOnMap; 