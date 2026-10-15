// The home page's globe: the Earth at night, turning slowly, with a glowing
// dot for every sighting. Drag to turn it. The positions come from the
// #globe-points JSON the page writes out, as [lat, lng] pairs of numbers.

import { mount } from '/assets/js/threejs.js';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';

// NASA's "Earth at night", from three.js's own examples. If it can't load,
// the globe stays plain night blue and the dots still show.
const LIGHTS = 'https://cdn.jsdelivr.net/gh/mrdoob/three.js@r186/examples/textures/planets/earth_lights_2048.png';
const RADIUS = 1.6;

function readPoints() {
	const el = document.getElementById('globe-points');
	try {
		const data = (el && JSON.parse(el.textContent)) || [];
		return data.filter((p) => Array.isArray(p) && Number.isFinite(p[0]) && Number.isFinite(p[1]));
	} catch (err) {
		return [];
	}
}

// Matches three.js's sphere texture mapping, so dots land on the right land.
function toVector(THREE, lat, lng, r) {
	const phi = THREE.MathUtils.degToRad(lat);
	const lambda = THREE.MathUtils.degToRad(lng);
	return new THREE.Vector3(
		r * Math.cos(phi) * Math.cos(lambda),
		r * Math.sin(phi),
		-r * Math.cos(phi) * Math.sin(lambda),
	);
}

// A soft round dot, drawn once on a canvas, for every point to share.
function glowTexture(THREE) {
	const size = 64;
	const canvas = document.createElement('canvas');
	canvas.width = canvas.height = size;
	const ctx = canvas.getContext('2d');
	const g = ctx.createRadialGradient(size / 2, size / 2, 0, size / 2, size / 2, size / 2);
	g.addColorStop(0, 'rgba(255,255,255,1)');
	g.addColorStop(0.25, 'rgba(255,255,255,0.8)');
	g.addColorStop(1, 'rgba(255,255,255,0)');
	ctx.fillStyle = g;
	ctx.fillRect(0, 0, size, size);
	return new THREE.CanvasTexture(canvas);
}

mount('globe', function ({ THREE, scene, camera, renderer, canvas, colour, still }) {
	const night = colour('--three-globe', '#16233a');
	const amber = colour('--three-primary', '#F5B841');
	const teal = colour('--three-accent', '#3FB8AF');

	camera.position.set(0, 0.6, 5.2);

	const globe = new THREE.Group();
	// Start with Africa facing us, tipped slightly towards the viewer.
	globe.rotation.y = -Math.PI / 2 - THREE.MathUtils.degToRad(20);
	globe.rotation.x = 0.15;
	scene.add(globe);

	const earthMaterial = new THREE.MeshBasicMaterial({ color: night });
	const earth = new THREE.Mesh(new THREE.SphereGeometry(RADIUS, 64, 48), earthMaterial);
	globe.add(earth);

	new THREE.TextureLoader().load(LIGHTS, function (texture) {
		texture.colorSpace = THREE.SRGBColorSpace;
		earthMaterial.map = texture;
		// Tinted, so the city lights read as a faint amber on night blue.
		earthMaterial.color = new THREE.Color(0xffffff).lerp(night, 0.35);
		earthMaterial.needsUpdate = true;
		if (still) renderer.render(scene, camera);
	});

	// Lines of latitude and longitude, very faint.
	globe.add(new THREE.Mesh(
		new THREE.SphereGeometry(RADIUS * 1.002, 36, 18),
		new THREE.MeshBasicMaterial({ color: teal, wireframe: true, transparent: true, opacity: 0.06 }),
	));

	// A halo round the edge, brighter where the globe meets space.
	globe.add(new THREE.Mesh(
		new THREE.SphereGeometry(RADIUS * 1.18, 64, 48),
		new THREE.ShaderMaterial({
			uniforms: { glow: { value: teal } },
			vertexShader: `
				varying vec3 vNormal;
				void main() {
					vNormal = normalize(normalMatrix * normal);
					gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
				}`,
			fragmentShader: `
				uniform vec3 glow;
				varying vec3 vNormal;
				void main() {
					float rim = pow(0.65 - dot(vNormal, vec3(0.0, 0.0, 1.0)), 3.0);
					gl_FragColor = vec4(glow, 1.0) * rim;
				}`,
			side: THREE.BackSide,
			blending: THREE.AdditiveBlending,
			transparent: true,
			depthWrite: false,
		}),
	));

	// Every sighting as one Points object, never a mesh each.
	const data = readPoints();
	const positions = new Float32Array(data.length * 3);
	data.forEach(function (p, i) {
		toVector(THREE, p[0], p[1], RADIUS * 1.01).toArray(positions, i * 3);
	});
	const geometry = new THREE.BufferGeometry();
	geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
	const dots = new THREE.PointsMaterial({
		color: amber,
		map: glowTexture(THREE),
		size: 0.11,
		sizeAttenuation: true,
		transparent: true,
		depthWrite: false,
		blending: THREE.AdditiveBlending,
	});
	globe.add(new THREE.Points(geometry, dots));

	// A few hundred stars behind it, for the night sky.
	const starCount = 600;
	const stars = new Float32Array(starCount * 3);
	for (let i = 0; i < starCount; i++) {
		const v = new THREE.Vector3().randomDirection().multiplyScalar(14 + Math.random() * 10);
		v.toArray(stars, i * 3);
	}
	const starGeometry = new THREE.BufferGeometry();
	starGeometry.setAttribute('position', new THREE.BufferAttribute(stars, 3));
	const starField = new THREE.Points(starGeometry, new THREE.PointsMaterial({
		color: 0xf2e9d8,
		size: 0.05,
		transparent: true,
		opacity: 0.7,
		depthWrite: false,
	}));
	scene.add(starField);

	const controls = new OrbitControls(camera, canvas);
	controls.enableZoom = false;
	controls.enablePan = false;
	controls.enableDamping = !still;
	controls.rotateSpeed = 0.5;
	controls.autoRotate = !still;
	controls.autoRotateSpeed = 0.6;
	controls.minPolarAngle = Math.PI * 0.2;
	controls.maxPolarAngle = Math.PI * 0.8;
	controls.addEventListener('change', function () {
		if (still) renderer.render(scene, camera);
	});

	return function (time) {
		// The dots breathe a little, as though they were glowing.
		dots.size = 0.11 + Math.sin(time * 2) * 0.015;
		starField.rotation.y = time * 0.005;
		controls.update();
	};
});
