// Flash Player Version Detection - Rev 1.5
// Detect Client Browser type
// Copyright(c) 2005-2006 Adobe Macromedia Software, LLC. All rights reserved.
// embedded inside a closure by N. Landsteiner, www.masswerk.at
// Web AudioManager (c) N.Landsteiner, wwww.masswerk.at
var AC_OETags = new function() {

var isIE  = (navigator.appVersion.indexOf("MSIE") != -1) ? true : false;
var isWin = (navigator.appVersion.toLowerCase().indexOf("win") != -1) ? true : false;
var isOpera = (navigator.userAgent.indexOf("Opera") != -1) ? true : false;

function ControlVersion()
{
	var version;
	var axo;
	var e;

	// NOTE : new ActiveXObject(strFoo) throws an exception if strFoo isn't in the registry

	try {
		// version will be set for 7.X or greater players
		axo = new ActiveXObject("ShockwaveFlash.ShockwaveFlash.7");
		version = axo.GetVariable("$version");
	} catch (e) {
	}

	if (!version)
	{
		try {
			// version will be set for 6.X players only
			axo = new ActiveXObject("ShockwaveFlash.ShockwaveFlash.6");
			
			// installed player is some revision of 6.0
			// GetVariable("$version") crashes for versions 6.0.22 through 6.0.29,
			// so we have to be careful. 
			
			// default to the first public version
			version = "WIN 6,0,21,0";

			// throws if AllowScripAccess does not exist (introduced in 6.0r47)		
			axo.AllowScriptAccess = "always";

			// safe to call for 6.0r47 or greater
			version = axo.GetVariable("$version");

		} catch (e) {
		}
	}

	if (!version)
	{
		try {
			// version will be set for 4.X or 5.X player
			axo = new ActiveXObject("ShockwaveFlash.ShockwaveFlash.3");
			version = axo.GetVariable("$version");
		} catch (e) {
		}
	}

	if (!version)
	{
		try {
			// version will be set for 3.X player
			axo = new ActiveXObject("ShockwaveFlash.ShockwaveFlash.3");
			version = "WIN 3,0,18,0";
		} catch (e) {
		}
	}

	if (!version)
	{
		try {
			// version will be set for 2.X player
			axo = new ActiveXObject("ShockwaveFlash.ShockwaveFlash");
			version = "WIN 2,0,0,11";
		} catch (e) {
			version = -1;
		}
	}
	
	return version;
}

// JavaScript helper required to detect Flash Player PlugIn version information
function GetSwfVer(){
	// NS/Opera version >= 3 check for Flash plugin in plugin array
	var flashVer = -1;
	
	if (navigator.plugins != null && navigator.plugins.length > 0) {
		if (navigator.plugins["Shockwave Flash 2.0"] || navigator.plugins["Shockwave Flash"]) {
			var swVer2 = navigator.plugins["Shockwave Flash 2.0"] ? " 2.0" : "";
			var flashDescription = navigator.plugins["Shockwave Flash" + swVer2].description;			
			var descArray = flashDescription.split(" ");
			var tempArrayMajor = descArray[2].split(".");
			var versionMajor = tempArrayMajor[0];
			var versionMinor = tempArrayMajor[1];
			if ( descArray[3] != "" ) {
				tempArrayMinor = descArray[3].split("r");
			} else {
				tempArrayMinor = descArray[4].split("r");
			}
			var versionRevision = tempArrayMinor[1] > 0 ? tempArrayMinor[1] : 0;
			var flashVer = versionMajor + "." + versionMinor + "." + versionRevision;
		}
	}
	// MSN/WebTV 2.6 supports Flash 4
	else if (navigator.userAgent.toLowerCase().indexOf("webtv/2.6") != -1) flashVer = 4;
	// WebTV 2.5 supports Flash 3
	else if (navigator.userAgent.toLowerCase().indexOf("webtv/2.5") != -1) flashVer = 3;
	// older WebTV supports Flash 2
	else if (navigator.userAgent.toLowerCase().indexOf("webtv") != -1) flashVer = 2;
	else if ( isIE && isWin && !isOpera ) {
		flashVer = ControlVersion();
	}	
	return flashVer;
}

// When called with reqMajorVer, reqMinorVer, reqRevision returns true if that version or greater is available
function DetectFlashVer(reqMajorVer, reqMinorVer, reqRevision)
{
	versionStr = GetSwfVer();
	if (versionStr == -1 ) {
		return false;
	} else if (versionStr != 0) {
		if(isIE && isWin && !isOpera) {
			// Given "WIN 2,0,0,11"
			tempArray         = versionStr.split(" "); 	// ["WIN", "2,0,0,11"]
			tempString        = tempArray[1];			// "2,0,0,11"
			versionArray      = tempString.split(",");	// ['2', '0', '0', '11']
		} else {
			versionArray      = versionStr.split(".");
		}
		var versionMajor      = versionArray[0];
		var versionMinor      = versionArray[1];
		var versionRevision   = versionArray[2];

        	// is the major.revision >= requested major.revision AND the minor version >= requested minor
		if (versionMajor > parseFloat(reqMajorVer)) {
			return true;
		} else if (versionMajor == parseFloat(reqMajorVer)) {
			if (versionMinor > parseFloat(reqMinorVer))
				return true;
			else if (versionMinor == parseFloat(reqMinorVer)) {
				if (versionRevision >= parseFloat(reqRevision))
					return true;
			}
		}
		return false;
	}
}

function AC_AddExtension(src, ext)
{
  if (src.indexOf('?') != -1)
    return src.replace(/\?/, ext+'?'); 
  else
    return src + ext;
}

function AC_Generateobj(objAttrs, params, embedAttrs) 
{ 
    var str = '';
    if (isIE && isWin && !isOpera)
    {
  		str += '<object ';
  		for (var i in objAttrs)
  			str += i + '="' + objAttrs[i] + '" ';
  		for (var i in params)
  			str += '><param name="' + i + '" value="' + params[i] + '" /> ';
  		str += '></object>';
    } else {
  		str += '<embed ';
  		for (var i in embedAttrs)
  			str += i + '="' + embedAttrs[i] + '" ';
  		str += '> </embed>';
    }

    document.write(str);
}

function AC_FL_RunContent(){
  var ret = 
    AC_GetArgs
    (  arguments, ".swf", "movie", "clsid:d27cdb6e-ae6d-11cf-96b8-444553540000"
     , "application/x-shockwave-flash"
    );
  AC_Generateobj(ret.objAttrs, ret.params, ret.embedAttrs);
}

function AC_GetArgs(args, ext, srcParamName, classid, mimeType){
  var ret = new Object();
  ret.embedAttrs = new Object();
  ret.params = new Object();
  ret.objAttrs = new Object();
  for (var i=0; i < args.length; i=i+2){
    var currArg = args[i].toLowerCase();    

    switch (currArg){	
      case "classid":
        break;
      case "pluginspage":
        ret.embedAttrs[args[i]] = args[i+1];
        break;
      case "src":
      case "movie":	
        args[i+1] = AC_AddExtension(args[i+1], ext);
        ret.embedAttrs["src"] = args[i+1];
        ret.params[srcParamName] = args[i+1];
        break;
      case "onafterupdate":
      case "onbeforeupdate":
      case "onblur":
      case "oncellchange":
      case "onclick":
      case "ondblClick":
      case "ondrag":
      case "ondragend":
      case "ondragenter":
      case "ondragleave":
      case "ondragover":
      case "ondrop":
      case "onfinish":
      case "onfocus":
      case "onhelp":
      case "onmousedown":
      case "onmouseup":
      case "onmouseover":
      case "onmousemove":
      case "onmouseout":
      case "onkeypress":
      case "onkeydown":
      case "onkeyup":
      case "onload":
      case "onlosecapture":
      case "onpropertychange":
      case "onreadystatechange":
      case "onrowsdelete":
      case "onrowenter":
      case "onrowexit":
      case "onrowsinserted":
      case "onstart":
      case "onscroll":
      case "onbeforeeditfocus":
      case "onactivate":
      case "onbeforedeactivate":
      case "ondeactivate":
      case "type":
      case "codebase":
        ret.objAttrs[args[i]] = args[i+1];
        break;
      case "id":
      case "width":
      case "height":
      case "align":
      case "vspace": 
      case "hspace":
      case "class":
      case "title":
      case "accesskey":
      case "name":
      case "tabindex":
        ret.embedAttrs[args[i]] = ret.objAttrs[args[i]] = args[i+1];
        break;
      default:
        ret.embedAttrs[args[i]] = ret.params[args[i]] = args[i+1];
    }
  }
  ret.objAttrs["classid"] = classid;
  if (mimeType) ret.embedAttrs["type"] = mimeType;
  return ret;
}

	this.ControlVersion = ControlVersion;
	this.GetSwfVer = GetSwfVer;
	this.DetectFlashVer = DetectFlashVer;
	this.AC_AddExtension = AC_AddExtension;
	this.AC_Generateobj = AC_Generateobj;
	this.AC_FL_RunContent = AC_FL_RunContent;
	this.AC_GetArgs = AC_GetArgs;
};

/* web audio */

var AudioManager = (function() {

function Sound(name, type, baseVolume,volume, loops, poolSize, delay) {
	this.audios=new Array();
	this.name=name;
	this.type=type;
	this.baseVolume=baseVolume;
	this.volume=volume;
	this.loops=loops || false;
	this.poolSize=poolSize || 1;
	this.delay=delay || 0;
	this.poolCnt=0;
	this.playing=false;
	this.timer=0;
	this.lastPlayed=0;
	this.buffer=null;
	this.sound=null;
	this.load();
}
Sound.prototype= {
	startFromCallback: false,
	ctx: null,
	masterGain: null,
	createMasterGain: function() {
		var context=Sound.prototype.ctx;
		if (!context) context=Sound.prototype.ctx=new AudioContext();
		if (context.createGain) {
			Sound.prototype.masterGain=context.createGain();
		}
		else {
			Sound.prototype.masterGain=context.createGainNode();
		}
		Sound.prototype.masterGain.connect(context.destination);
	},
	setMasterGain: function(v) {
		if (!Sound.prototype.masterGain) Sound.prototype.createMasterGain();
		Sound.prototype.masterGain.gain.value=v;
	},
	load: function(autoplay) {
		if (window.AudioContext) {
			try {
				if (!Sound.prototype.masterGain) {
					Sound.prototype.createMasterGain();
					Sound.prototype.setMasterGain(this.volume);
				}
				var context=Sound.prototype.ctx;
				var request = new XMLHttpRequest();
				request.open('GET', soundPath+this.name+this.type, true);
				request.responseType = 'arraybuffer';
				var ref=this;
				request.onload = (this.startFromCallback)?
					function() {
						ref.buffer=request.response;
					} :
					function() {
						context.decodeAudioData(request.response, function(buffer) {
							ref.buffer = buffer;
						}, ref.log);
					};
				request.send();
			}
			catch (e) {
				this.log(e);
			}
		}
	},
	play: function(ignoreDelay) {
		if (window.AudioContext) {
			if (!this.buffer) return;
			try {
				if (this.volume==0 || this.playing) return;
				if (this.delay && !ignoreDelay) {
					var t=new Date().getTime();
					var d=(this.lastPlayed)? this.delay-(t-this.lastPlayed):0;
					if (d>0) {
						if (!this.timer) {
							var that=this;
							this.timer=setTimeout(function() { that.play(); }, d);
						}
						return;
					}
					this.lastPlayed=t;
				}
				if (this.timer) {
					clearTimeout(this.timer);
					this.timer=0;
				}
				var context=Sound.prototype.ctx;
				var source = context.createBufferSource();
				if (this.loops) {
					source.loop = true;
					this.playing = true;
				}
				var gainNode;
				if (context.createGain) {
					gainNode=context.createGain();
				}
				else {
					gainNode=context.createGainNode();
				}
				source.connect(gainNode);
				gainNode.connect(Sound.prototype.masterGain);
				gainNode.gain.value=this.baseVolume;
				if (this.startFromCallback) {
					var ref=this;
					context.decodeAudioData(this.buffer, function(buffer) {
						source.buffer = buffer;
						if (source.start) {
							source.start(0);
						}
						else {
							source.noteOn(0);
						}
					}, ref.log);
				}
				else {
					source.buffer = this.buffer;
					if (source.start) {
						source.start(0);
					}
					else {
						source.noteOn(0);
					}
				}
				this.sound=source;
			}
			catch (e) {
				this.log(e);
			}
		}
	},
	stop: function() {
		if (window.AudioContext) {
			try {
				if (this.sound) {
					if (this.sound.stop) {
						this.sound.stop(0);
					}
					else {
						this.sound.noteOff(0);
					}
					this.sound=null;
				}
			}
			catch (e) {
				this.log(e);
			}
			this.playing=false;
			if (this.timer) {
				clearTimeout(this.timer);
				this.timer=0;
			}
		}
	},
	log: function (e) {
		if (window.console) console.log('Audio Exception: ',e);
	}
};
	var type, sounds=new Object(), volume=3/7, hasWebAudio=false,
		soundData={}, soundPath='';
	function setup() {
		if (!window.AudioContext) {
			var vendors=['webkit', 'moz', 'o', 'ms'];
			for (var i=0; i<vendors.length; i++) {
				var api=window[vendors[i]+'AudioContext'];
				if (api) {
					window.AudioContext=api;
					break;
				}
			}
		}
		hasWebAudio=Boolean(window.AudioContext);
		if (!hasWebAudio) return;
		if (hasWebAudio && typeof navigator !== 'undefined' && navigator.userAgent) {
			if (navigator.userAgent.indexOf('Chrome')!==-1) {
				if (parseInt(navigator.userAgent.replace(/^.*?\bChrome\/([0-9]+).*$/, '$1'),10)>=32) {
					Sound.prototype.startFromCallback=true;
				}
			}
			else if (navigator.userAgent.match(/\bVersion\/[0-9]+\.[0-9\.]+ (Mobile\/\w+ )?Safari\//)) {
				if (parseInt(navigator.userAgent.replace(/^.*?\bVersion\/([0-9]+).*$/, '$1'),10)>=9) {
					Sound.prototype.startFromCallback=true;
				}
			}
		}
		var a=document.createElement('audio');
		if (!a || !a.canPlayType) {
			if (hasWebAudio) type='.wav';
		}
		else if (a.canPlayType('audio/mpeg')!='') {
			type='.mp3';
		}
		else if (a.canPlayType('audio/ogg')!='') {
			type='.ogg';
		}
		else if (a.canPlayType('audio/wav')!='' || a.canPlayType('audio/x-wav')!='') {
			type='.wav';
		}
	}
	function load(sname) {
		if (!type) return;
		if (!sounds[sname]) {
			var sdata= soundData[sname];
			if (sdata) {
				sounds[sname]=new Sound(sdata.name, type, sdata.baseVolume, volume, sdata.loops, sdata.poolSize, sdata.delay);
			}
		}
	}
	function play(sname, ignoreDelay) {
		if (!type) return;
		if (sounds[sname]) {
			sounds[sname].play(ignoreDelay);
		}
	}
	function setVolume(v, isfinal) {
		if (!type) return;
		if (v!=volume) {
			if (hasWebAudio) Sound.prototype.setMasterGain(v);
			volume=v;
		}
	}
	function stopSound(sname) {
		if (!type) return;
		if (sounds[sname]) sounds[sname].stop();
	}
	function stopAllSounds() {
		if (!type) return;
		for (var i in sounds) sounds[i].stop();
	}
	function destroy() {
		stopAllSounds();
		sounds=null;
	}
	function preload(path, data) {
		if (hasWebAudio && typeof data==='object') {
			if (path) {
				if (!path.match(/\/$/)) path+='/';
				soundPath=path;
			}
			soundData=data;
			for (var i in data) {
				var d=data[i];
				if (typeof d==='object' && d.name) load(d.name);
			}
		}
	}
	function supported() {
		return Boolean(hasWebAudio && type);
	}
	setup();
	return {
		'load': preload,
		'play': play,
		'stopSound': stopSound,
		'stopAllSounds': stopAllSounds,
		'destroy': destroy,
		'setVolume': setVolume,
		'supported': supported
	};
})();
