// Index page JavaScript
async function loadHeroStats() {
  try {
    const response = await fetch('/api/stats');
    const data = await response.json();
    
    if (data.success) {
      // Update hero stats
      const totalReportsElement = document.getElementById('total-reports');
      const activeCitizensElement = document.getElementById('active-citizens');
      
      if (totalReportsElement) {
        totalReportsElement.textContent = data.total_reports.toLocaleString('pt-BR');
      }
      
      if (activeCitizensElement) {
        activeCitizensElement.textContent = data.active_citizens.toLocaleString('pt-BR');
      }
    }
  } catch (error) {
    console.error('Error loading hero stats:', error);
  }
}

// Load stats when the page loads
document.addEventListener('DOMContentLoaded', loadHeroStats);
