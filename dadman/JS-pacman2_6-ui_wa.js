/*
  JavaScript-PacMan 2.6 UI-Controls
  (c) N.Landsteiner 2007-2012; www.masswerk.at
  All rights reserved.
*/

// bridge & controls to flash sound

var JSPacManUI = new function() {

var mazeOverlayId = 'pacmanMazeoverlay';
var soundSliderId='pacmanSndslider';
var soundSliderBgId='pacmanSndsliderbg';

var soundsAvailable=false;
var soundApp=null;
var useSound=false;
var useWebAudio=false;

if (AudioManager.supported()) {
	useWebAudio=soundsAvailable=true;
	AudioManager.load('jspacsounds/', {
		'bonuslife': { 'name': 'bonuslife', 'baseVolume': 1, 'loops': false },
		'doteat': { 'name': 'doteat', 'baseVolume': 0.95, 'loops': false, 'poolSize': 2 },
		'endlevel': { 'name': 'endlevel', 'baseVolume': 1, 'loops': false },
		'ghostcatch': { 'name': 'ghostcatch', 'baseVolume': 1, 'loops': false, 'poolSize': 4 },
		'pacend': { 'name': 'pacend', 'baseVolume': 1, 'loops': false },
		'pillbeam': { 'name': 'pillbeam', 'baseVolume': 1, 'loops': false }
	});
	document.addEventListener('DOMContentLoaded', initSnd, false);
}

function setSoundsAvailable() {
	soundApp=getSwfApp('jspacsounds');
	if (soundApp) {
		soundsAvailable=true;
		initSnd();
	}
}

function getSwfApp(appName) {
	if (document.getElementById) {
		return document.getElementById(appName);
	}
	if (navigator.appName.indexOf ("Microsoft") !=-1) {
		return window[appName];
	}
	else {
		return document[appName];
	}
}

function playSound(s) {
	if (useWebAudio) {
		AudioManager.play(s);
	}
	else if (soundsAvailable) {
		try {
			soundApp.playSound(s);
		}
		catch (e) {}
	}
}

function setVolume(v) {
	if (useWebAudio) {
		AudioManager.setVolume(v);
	}
	else if (soundsAvailable) {
		try {
			soundApp.setVolume(v);
		}
		catch (e) {}
	}
}

var sndIntialVolume=.25;
var sndSliderZero=1;
var sndSliderMin=11;
var sndSliderMax=81;
var sndSliderMinValue=0.05;
var sndSliderMaxValue=0.9;

var sndSliderSpan=sndSliderMax-sndSliderMin;
var sndSliderValueSpan=sndSliderMaxValue-sndSliderMinValue;

var sliderElement=null;
var sliderDragging=false;
var sliderX;
var sliderDragMouseX=0;
var sliderTimer=null;

function initSnd() {
	var sliderbg;
	if (document.getElementById) {
		sliderElement=document.getElementById(soundSliderId);
		sliderbg=document.getElementById(soundSliderBgId);
	}
	else if (document.all) {
		sliderElement=document.all[soundSliderId];
		sliderbg=document.all[soundSliderBgId];
	}
	if (sliderElement) {
		setVolume(sndIntialVolume);
		if (sndIntialVolume>=sndSliderMinValue) {
			sliderX= sndIntialVolume/sndSliderValueSpan * sndSliderSpan + sndSliderMin;
			useSound=true;
			if (JSPacMan) JSPacMan.useSound(true);
		}
		else {
			sliderX=sndSliderZero;
			useSound=false;
			if (JSPacMan) JSPacMan.useSound(false);
		}
		sliderElement.style.left=sliderX+'px';
		registerEvent(sliderElement, 'mousedown', sliderDown);
		if (sliderbg) registerEvent(sliderbg, 'mousedown', sliderBgDown);
	}
}

function sliderDown(e) {
	if (window.event) e = window.event;
	sliderDragMouseX = (typeof e.clientX != 'undefined')? e.clientX:e.pageX;
	registerEvent(document, 'mousemove', sliderDrag);
	registerEvent(document, 'mouseup', sliderUp);
	sliderDragging=true;
	if (window.setOpacity) setOpacity(sliderElement, 0.85);
	stopEvent(e);
	return false;
}

function sliderDrag(e) {
	if (window.event) e = window.event;
	if (sliderDragging) {
		var x=(typeof e.clientX != 'undefined')? e.clientX:e.pageX;
		var dx=x-sliderDragMouseX;
		var sx=sliderX;
		setVolumeSlider(sliderX+dx, false);
		sliderDragMouseX+=sliderX-sx;
	}
	stopEvent(e);
	return false;
}

function sliderUp(e) {
	if (window.event) e = window.event;
	if (sliderDragging) {
		var x=(typeof e.clientX != 'undefined')? e.clientX:e.pageX;
		var dx=x-sliderDragMouseX;
		setVolumeSlider(sliderX+dx, true);
		sliderDragging=false;
		releaseEvent(document, 'mousemove', sliderDrag);
		releaseEvent(document, 'mouseup', sliderUp);
		if (window.setOpacity) setOpacity(sliderElement, 1);
	}
	stopEvent(e);
	return false;
}

function setVolumeSlider(x, isFinal) {
	if (x>sndSliderMax) {
		x=sndSliderMax;
	}
	else if (x<sndSliderZero) {
		x=sndSliderZero;
	}
	if (x<sndSliderMin) {
		if (useSound) setVolume(0);
		useSound=false;
		if (JSPacMan) JSPacMan.useSound(false);
		if (isFinal) x=sndSliderZero;
	}
	else {
		useSound=true;
		if (JSPacMan) JSPacMan.useSound(true);
		if (x!=sliderX) {
			setVolume( (x-sndSliderMin)/sndSliderSpan*sndSliderValueSpan+sndSliderMinValue );
		}
	}
	sliderElement.style.left=x+'px';
	sliderX=x;
}

function sliderBgDown(e) {
	if (sliderDragging) return;
	if (window.event) e = window.event;
	var x;
	if (typeof e.offsetX != 'undefined') {
		x=e.offsetX;
	}
	else if (typeof e.layerX != 'undefined') {
		x=e.layerX;
	}
	else return;
	x-=4;
	if (x<sndSliderZero) {
		x=sndSliderZero;
	}
	else if (x>sndSliderMax) {
		x=sndSliderMax;
	}
	if (sliderTimer) clearTimeout(sliderTimer);
	sliderMoveTo(x);
	stopEvent(e);
	return false;
}

function sliderMoveTo(x) {
	var dx=x-sliderX;
	if (Math.abs(dx)>3) {
		sliderX+=(dx>0)? 3:-3;
		sliderElement.style.left=sliderX+'px';
		sliderTimer=setTimeout( function() { sliderMoveTo(x); }, 10);
	}
	else {
		setVolumeSlider(x, true);
	}
}

function registerEvent(obj, eventType, handler) {
	if (obj.addEventListener) {
		obj.addEventListener(eventType.toLowerCase(), handler, false);
	}
	else if (obj.attachEvent) {
		obj.attachEvent('on'+eventType.toLowerCase(), handler);
	}
	else {
		var et=eventType.toUpperCase();
		if ((obj.Event) && (obj.Event[et]) && (obj.captureEvents)) obj.captureEvents(Event[et]);
		obj['on'+eventType.toLowerCase()]=handler;
	}
}
function releaseEvent(obj, eventType, handler) {
	if (obj.removeEventListener) {
		obj.removeEventListener(eventType.toLowerCase(), handler, false);
	}
	else if (obj.detachEvent) {
		obj.detachEvent('on'+eventType.toLowerCase(), handler);
	}
	else {
		var et=eventType.toUpperCase();
		if ((obj.Event) && (obj.Event[et]) && (obj.releaseEvents)) obj.releaseEvents(Event[et]);
		et='on'+eventType.toLowerCase();
		if ((obj[et]) && (obj[et]==handler)) obj.et=null;
	}
}
function stopEvent(e) {
	if (e.preventDefault) e.preventDefault();
	if (e.stopPropagation) e.stopPropagation();
	e.cancelBubble=true;
	e.returnValue=false;
}

// touch srceen controls

var touchX=0;
var touchY=0;
var touchTime=0;
var touchDown=false;
var touchLastX=0;
var touchLastY=0;

function activateTouchControls() {
	var touchtarget=document.getElementById(mazeOverlayId);
	if (touchtarget) {
		if (RegExp(/\bmobile\b/i).test(navigator.userAgent)) {
			registerEvent(touchtarget, 'touchstart', mobileTouchHandlerDown);
			registerEvent(document, 'touchmove', mobileTouchHandlerMove);
			registerEvent(document, 'touchend', mobileTouchHandlerUp);
			registerEvent(document, 'touchcancel', mobileTouchHandlerCancel);
		}
		else {
			registerEvent(touchtarget, 'mousedown', touchHandlerDown);
			registerEvent(document, 'mousemove', touchHandlerMove);
			registerEvent(document, 'mouseup', touchHandlerUp);
		}
	}
}

function touchProcessDelta(dx, dy, minDelta) {
	var adx=Math.abs(dx);
	var ady=Math.abs(dy);
	if (adx>ady) {
		if (adx>minDelta) JSPacMan.setPacDirection((dx>0)? 1:2);
	}
	else {
		if (ady>minDelta) JSPacMan.setPacDirection((dy>0)? 8:4);
	}
}

function touchProcessDoubleClick(dx, dy, maxDelta, maxTicks) {
	if (Math.abs(dx)<=maxDelta && Math.abs(dy)<=maxDelta) {
		var now=new Date();
		var t=now.getTime();
		if (t-touchTime<maxTicks) {
			JSPacMan.pause();
			touchTime=0;
		}
		else {
			touchTime=t;
		}
	}
	else {
		touchTime=0;
	}
}

function touchHandlerDown(e) {
	if (window.event) e=window.event;
	if (e.clientX != undefined) {
		touchX=e.clientX;
		touchY=e.clientY;
	}
	else if (e.pageX != undefined) {
		touchX=e.pageX;
		touchY=e.pageY;
	}
	else {
		stopEvent(e);
		return false;
	}
	touchDown=true;
	stopEvent(e);
	return false;
}

function touchHandlerMove(e) {
	if (!touchDown) return;
	if (window.event) e=window.event;
	if (e.clientX != undefined) {
		touchLastX=e.clientX;
		touchLastY=e.clientY;
	}
	else if (e.pageX != undefined) {
		touchLastX=e.pageX;
		touchLastY=e.pageY;
	}
	else {
		stopEvent(e);
		return false;
	}
	touchProcessDelta(touchLastX-touchX, touchLastY-touchY, 3);
	stopEvent(e);
	return false;
}

function touchHandlerUp(e) {
	if (!touchDown) return;
	if (window.event) e=window.event;
	if (e.clientX != undefined) {
		touchLastX=e.clientX;
		touchLastY=e.clientY;
	}
	else if (e.pageX != undefined) {
		touchLastX=e.pageX;
		touchLastY=e.pageY;
	}
	else {
		stopEvent(e);
		return false;
	}
	touchProcessDoubleClick(touchLastX-touchX, touchLastY-touchY, 3, 600);
	touchDown=false;
	stopEvent(e);
	return false;
}

function mobileTouchHandlerDown(e) {
	if (e && e.touches && e.touches.length) {
		if (e.touches.length > 1) return true;
		var t=e.touches[0];
		if (t.clientX != undefined) {
			touchX=t.clientX;
			touchY=t.clientY;
		}
		else if (t.pageX != undefined) {
			touchX=t.pageX;
			touchY=t.pageY;
		}
		else {
			stopEvent(e);
			return false;
		}
		touchDown=true;
		stopEvent(e);
		return false;
	}
}

function mobileTouchHandlerMove(e) {
	if (!touchDown) return;
	var t=e.touches[0];
	if (t.clientX != undefined) {
		touchLastX=t.clientX;
		touchLastY=t.clientY;
	}
	else if (t.pageX != undefined) {
		touchLastX=t.pageX;
		touchLastY=t.pageY;
	}
	else {
		stopEvent(e);
		return false;
	}
	touchProcessDelta(touchLastX-touchX, touchLastY-touchY, 10);
	stopEvent(e);
	return false;
}

function mobileTouchHandlerUp(e) {
	if (!touchDown) return;
	var t;
	if (e.touches && e.touches.length) {
		t=e.touches[0];
	}
	else if (e.changedTouches && e.changedTouches.length) {
		t=e.changedTouches[0];
	}
	if (t) {
		if (t.clientX != undefined) {
			touchLastX=t.clientX;
			touchLastY=t.clientY;
		}
		else if (t.pageX != undefined) {
			touchLastX=t.pageX;
			touchLastY=t.pageY;
		}
		else {
			stopEvent(e);
			return false;
		}
		touchProcessDoubleClick(touchLastX-touchX, touchLastY-touchY, 10, 600);
	}
	touchDown=false;
	stopEvent(e);
	return false;
}

function mobileTouchHandlerCancel(e) {
	if (touchDown) touchDown=false;
}

function saveLastTouch() {
	touchX=touchLastX;
	touchY=touchLastY;
}

// export external API

return {
	setSoundsAvailable: setSoundsAvailable,
	soundActive:function() { return useSound; },
	activateTouchControls: activateTouchControls,
	saveLastTouch: saveLastTouch,
	playSound: playSound,
	setVolume: setVolume
};

// end of closure
};

// Alias hand-shake API for Flash Player
var setSoundsAvailable = JSPacManUI.setSoundsAvailable;