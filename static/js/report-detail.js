// Report Detail Map functionality
let map;
let marker;
let AdvancedMarkerElement;
let PinElement;

// This function will be called by Google Maps API when it loads
async function initReportDetailMap() {
    // Get coordinates from the map element's data attributes
    const mapElement = document.getElementById('map');
    if (!mapElement) {
        console.error('Map element not found');
        return;
    }
    
    const latitude = parseFloat(mapElement.dataset.latitude);
    const longitude = parseFloat(mapElement.dataset.longitude);
    
    if (isNaN(latitude) || isNaN(longitude)) {
        console.error('Invalid coordinates:', latitude, longitude);
        return;
    }
    
    const location = { lat: latitude, lng: longitude };
    
    try {
        // Import the marker library
        const { AdvancedMarkerElement: AME, PinElement: PE } = await google.maps.importLibrary("marker");
        AdvancedMarkerElement = AME;
        PinElement = PE;
        
        // Create the map
        map = new google.maps.Map(document.getElementById('map'), {
            zoom: 15,
            center: location,
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
        
        // Create pin element for the marker
        const pinElement = new PinElement({
            background: '#dc3545', // Red color for report location
            borderColor: 'white',
            glyphColor: 'white',
            scale: 1.2
        });
        
        // Create the marker
        marker = new AdvancedMarkerElement({
            position: location,
            map: map,
            content: pinElement.element,
            title: 'Localização da Denúncia'
        });
        
    } catch (error) {
        console.error('Error creating map:', error);
    }
}

// Make function available globally
window.initReportDetailMap = initReportDetailMap;
