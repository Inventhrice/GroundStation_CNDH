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

var updateURI = "http://localhost:8080/update";
var update = new EventSource(updateventURI);
update.addEventListener("message", function(event) {
    const telemetry = JSON.parse(event.data);
    document.getElementById("x-coordinate").innerText(telemetry.coordinates.x);
    document.getElementById("y-coordinate").innerText(telemetry.coordinates.y);
    document.getElementById("z-coordinate").innerText(telemetry.coordinates.z);
    document.getElementById("pitch").innerText(telemetry.rotations.p);
    document.getElementById("yaw").innerText(telemetry.rotations.y);
    document.getElementById("roll").innerText(telemetry.roations.r);
    document.getElementById("temp").innerText(telemetry.temp);
    document.getElementById("payload-power").innerText(telemetry.status.payloadPower);
    document.getElementById("data-waiting").innerText(telemetry.status.dataWaiting);
    document.getElementById("charge-status").innerText(telemetry.status.chargeStatus);
    document.getElementById("voltage").innerText(telemetry.status.voltage);
});

window.addEventListener("beforeunload", function() {
    update.close();
});
