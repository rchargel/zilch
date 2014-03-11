var POS_X = 1800;
var POS_Y = 300;
var POS_Z = 1800;
var WIDTH = 1000;
var HEIGHT = 600;

var FOV = 45;
var NEAR = 1;
var FAR = 4000;

// some global variables and initialization code
// simple basic renderer
var renderer = new THREE.WebGLRenderer();
renderer.setSize(WIDTH,HEIGHT);
renderer.setClearColorHex(0x000000);

// add it to the target element
var mapDiv = document.getElementById("globe");
mapDiv.appendChild(renderer.domElement);

// setup a camera that points to the center
var camera = new THREE.PerspectiveCamera(FOV,WIDTH/HEIGHT,NEAR,FAR);
camera.position.set(POS_X,POS_Y, POS_Z);
camera.lookAt(new THREE.Vector3(0,0,0));

// create a basic scene and add the camera
var scene = new THREE.Scene();
scene.add(camera);

var light = new THREE.DirectionalLight(0x5555ff, 3.5, 500);

// we wait until the document is loaded before loading the
// density data.
$(document).ready(function()  {
	addLights();
	addEarth();
	addClouds();
	addDistribution();
	render();
});

function addEarth() {
	var spGeo = new THREE.SphereGeometry(600,50,50);
	var planetTexture = THREE.ImageUtils.loadTexture("/images/world-big-2-grey.jpg");
	var mat2 = new THREE.MeshPhongMaterial({
		map: planetTexture,
		shininess: 0.2
	});
	var sp = new THREE.Mesh(spGeo,mat2);
	scene.add(sp);
}

function addLights() {
	scene.add(light);
	light.position.set(POS_X,POS_Y,POS_Z);
}

function addDistribution() {
	var spGeo = new THREE.SphereGeometry(600,50,50);
	var cloudsTexture = THREE.ImageUtils.loadTexture("/map_transparent.png");
	var materialClouds = new THREE.MeshPhongMaterial( { color: 0xffffff, map: cloudsTexture, transparent:true, opacity: 1 } );

	meshClouds = new THREE.Mesh( spGeo, materialClouds );
	meshClouds.scale.set( 1.015, 1.015, 1.015 );
	scene.add( meshClouds );
}

function addClouds() {
	var spGeo = new THREE.SphereGeometry(600,50,50);
	var cloudsTexture = THREE.ImageUtils.loadTexture( "/images/earth_clouds_1024.png" );
	var materialClouds = new THREE.MeshPhongMaterial( { color: 0xffffff, map: cloudsTexture, transparent:true, opacity:0.3 } );

	meshClouds = new THREE.Mesh( spGeo, materialClouds );
	meshClouds.scale.set( 1.005, 1.005, 1.005 );
	scene.add( meshClouds );
}

function copy(obj) {
	var c = {};
	for (var k in obj) {
		c[k] = obj[k];
	}
	return c;
}

function render() {
	var timer = Date.now() * 0.0001;
	camera.position.x = (Math.cos(timer) * 1800);
	camera.position.z = (Math.sin(timer) * 1800);
	camera.lookAt(scene.position);
	light.position = copy(camera.position);
	light.position.y = camera.position.y - 300;
	light.lookAt(scene.position);
	renderer.render(scene, camera);
	requestAnimationFrame(render);
}
