var nav = document.getElementById('nav'),
    output = document.getElementById('output');

function Dir(dir) {
  this.name = dir.name;
  this.is_dir = dir.is_dir;
  this.is_markdown = dir.is_markdown;
  this.content = dir.content;
  
  for (i in dir.children) {
    var c = dir.children[i];
    
    if (c.is_dir || c.is_markdown) {
      var d = new Dir(c);
      nav.appendChild(d.element());
    }
  }
}

Dir.prototype.element = function() {
  var a = document.createElement('a');
  a.innerHTML = this.name;
  a.href = '#';
  a.className = 'mdl-navigation__link'
  a.onclick = function() {
    update(this);
    return false;
  }.bind(this);
  return a;
}

function newRequest(path, cb) {
  output.innerHTML = '<div id="p2" class="mdl-progress mdl-js-progress mdl-progress__indeterminate"></div>';
  var r = new XMLHttpRequest();
  r.addEventListener('load', cb)
  r.open('GET', '/-/' + path);
  r.send();
}

var parents = [];

function Update(path) {
  parents = path.split('/');
  update();
}

function update(dir) {
  nav.innerHTML = '';
  
  if (dir == '..') {
    parents.pop();
  }
  
  if (dir && dir != '..') {
    parents.push(dir.name);
  }

  history.pushState(null, null, "/" + parents.join('/'));

  newRequest(parents.join('/'), function() {
    try {
      var dir = JSON.parse(this.responseText);
    } catch (e) {
      console.error(e, this.responseText);
      return;
    }

    var d = new Dir(dir);

    if (d.is_markdown) {
      // TODO move to blackfriday.HTML.Code
      d.content = d.content.replace(/<pre><code class=/g, '<pre class="prettyprint"><code class="');
      output.innerHTML = d.content;
      PR.prettyPrint();
    } else {
      output.innerHTML = '';
    }
  })
}

function UpdateByHash() {
  if (window.location.hash)
    return;
  var hash = window.location.pathname;
  if (!hash)
    return false;
  Update(hash.substring(1));
  return true;
}

var KEYBOARD_T = 84;

window.onkeyup = function(e) {
  e.preventDefault();

  if (e.target.nodeName == 'INPUT') return;

  if (e.keyCode == KEYBOARD_T) {
    var path = prompt(">>>", parents.join('/'));
    if (path) {
      Update(path);  
    }
  }
}

window.onload = function() {
  if (!UpdateByHash()) {
    update();
  }
}

window.onpopstate = function(e) {
  UpdateByHash();
}