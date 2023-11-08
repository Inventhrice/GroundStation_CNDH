// Select container section to appear when clicked
var telemetryInputContainer = document.querySelector('.telemetry-input-container');
var coordinateInputContainer = document.querySelector('.coordinate-input-container');
var angleInputContainer = document.querySelector('.angle-input-container');

// Set buttons to be able to do actions
var setTelemetryBtn = document.getElementById('setTelemetryBtn');
var setCoordinateBtn = document.getElementById('setCoordinateBtn');
var setAngleBtn = document.getElementById('setAngleBtn');

// Don't display the containers by default
telemetryInputContainer.style.display = 'none';
coordinateInputContainer.style.display = 'none';
angleInputContainer.style.display = 'none';

// Add click event listener to the Change Telemetry button
setTelemetryBtn.addEventListener('click', function() {
    // Toggle the display property of the telemetry input container
    if (telemetryInputContainer.style.display === 'none' || telemetryInputContainer.style.display === '') {
        angleInputContainer.style.display = 'none';
        coordinateInputContainer.style.display = 'none';
        telemetryInputContainer.style.display = 'block';
    } else {
        telemetryInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Change Coordinate button
setCoordinateBtn.addEventListener('click', function() {
    // Toggle the display property of the coordinate input container
    if (coordinateInputContainer.style.display === 'none' || coordinateInputContainer.style.display === '') {
        angleInputContainer.style.display = 'none';
        telemetryInputContainer.style.display = 'none';
        coordinateInputContainer.style.display = 'block';
    } else {
        coordinateInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Change Angle button
setAngleBtn.addEventListener('click', function() {
    // Toggle the display property of the angle input container
    if (angleInputContainer.style.display === 'none' || angleInputContainer.style.display === '') {
        telemetryInputContainer.style.display = 'none';
        coordinateInputContainer.style.display = 'none';
        angleInputContainer.style.display = 'block';
    } else {
        angleInputContainer.style.display = 'none';
    }
});