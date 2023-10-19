// Select container section to appear when clicked
var telemetryInputContainer = document.querySelector('.telemetry-input-container');
var zoomInputContainer = document.querySelector('.zoom-input-container');
var tempInputContainer = document.querySelector('.temp-input-container');
var angleInputContainer = document.querySelector('.angle-input-container');

// Set buttons to be able to do actions
var changeTelemetryBtn = document.getElementById('changeTelemetryBtn');
var setZoomBtn = document.getElementById('setZoomBtn');
var setTempBtn = document.getElementById('setTempBtn');
var setAngleBtn = document.getElementById('setAngleBtn');

// Don't display the containers by default
telemetryInputContainer.style.display = 'none';
zoomInputContainer.style.display = 'none';
tempInputContainer.style.display = 'none';
angleInputContainer.style.display = 'none';

// Add click event listener to the Change Telemetry button
changeTelemetryBtn.addEventListener('click', function() {
    // Toggle the display property of the telemetry input container
    if (telemetryInputContainer.style.display === 'none' || telemetryInputContainer.style.display === '') {
        zoomInputContainer.style.display = 'none';
        tempInputContainer.style.display = 'none';
        angleInputContainer.style.display = 'none';
        telemetryInputContainer.style.display = 'block';
    } else {
        telemetryInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Set Zoom button
setZoomBtn.addEventListener('click', function() {
    // Toggle the display property of the zoom input container
    if (zoomInputContainer.style.display === 'none' || zoomInputContainer.style.display === '') {
        telemetryInputContainer.style.display = 'none';
        tempInputContainer.style.display = 'none';
        angleInputContainer.style.display = 'none';
        zoomInputContainer.style.display = 'block';
    } else {
        zoomInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Set Zoom button
setTempBtn.addEventListener('click', function() {
    // Toggle the display property of the zoom input container
    if (tempInputContainer.style.display === 'none' || tempInputContainer.style.display === '') {
        telemetryInputContainer.style.display = 'none';
        zoomInputContainer.style.display = 'none';
        angleInputContainer.style.display = 'none';
        tempInputContainer.style.display = 'block';
    } else {
        tempInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Change Telemetry button
setAngleBtn.addEventListener('click', function() {
    // Toggle the display property of the telemetry input container
    if (angleInputContainer.style.display === 'none' || angleInputContainer.style.display === '') {
        zoomInputContainer.style.display = 'none';
        tempInputContainer.style.display = 'none';
        telemetryInputContainer.style.display = 'none';
        angleInputContainer.style.display = 'block';
    } else {
        angleInputContainer.style.display = 'none';
    }
});