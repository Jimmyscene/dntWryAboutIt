import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import { Terminal } from 'xterm';
import * as  attach from 'xterm/lib/addons/attach/attach';

import 'xterm/src/xterm.css'

Terminal.applyAddon(attach);  // Apply the `attach` addon
const val = prompt("Debug?")
var term = new Terminal();
var socket = new WebSocket('ws://localhost:9000/pty?debug=' + val);

term.attach(socket);  // Attach the above socket to `term`

term.open(document.getElementById('xterm-container'));
ReactDOM.render(<App />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
