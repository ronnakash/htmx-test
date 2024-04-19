// elements
const levelInput = document.getElementById('levelInput');
const search = document.getElementById('search');
const filterForm = document.getElementById('filterForm');

// Function to update levelInput field
function updateLevel(level) {
  levelInput.value = level;
}

// trigger htmx event on search text change
search.addEventListener('keyup', function(event) {
  console.log("hi");
  event.preventDefault();
  const updateEvent = new Event("textUpdate");
  filterForm.dispatchEvent(updateEvent);
});
