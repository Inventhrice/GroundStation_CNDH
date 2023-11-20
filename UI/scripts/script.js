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
var update = new EventSource("/update");
update.addEventListener("message", function(event) {
    console.log("Update!");
    const telemetry = JSON.parse(event.data);
    document.getElementById("x-coordinate").innerHTML = telemetry.coordinates.x;
    document.getElementById("y-coordinate").innerHTML = telemetry.coordinates.y;
    document.getElementById("z-coordinate").innerHTML = telemetry.coordinates.z;
    document.getElementById("pitch").innerHTML = telemetry.rotations.p;
    document.getElementById("yaw").innerHTML = telemetry.rotations.y;
    document.getElementById("roll").innerHTML = telemetry.rotations.r;
    document.getElementById("temp").innerHTML = telemetry.temp;
    document.getElementById("payload-power").innerHTML = telemetry.status.payloadPower;
    document.getElementById("data-waiting").innerHTML = telemetry.status.dataWaiting;
    document.getElementById("charge-status").innerHTML = telemetry.status.chargeStatus;
    document.getElementById("voltage").innerHTML = telemetry.status.voltage;
});

window.addEventListener("beforeunload", function() {
    update.close();
});
