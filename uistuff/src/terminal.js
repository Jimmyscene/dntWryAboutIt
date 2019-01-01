import { XTerm, Terminal } from 'react-xterm'
import React, { Component } from 'react';
import * as  attach from './attach';
import * as fit from 'xterm/lib/addons/fit/fit';
import * as fullscreen from 'xterm/lib/addons/fullscreen/fullscreen';
import * as search from 'xterm/lib/addons/search/search';
import 'xterm/src/xterm.css'

Terminal.applyAddon(attach)
Terminal.applyAddon(fit)
Terminal.applyAddon(fullscreen)
Terminal.applyAddon(search)


class MyTerminal extends Component {
    componentDidUpdate(prevProps) {
        if (prevProps.debug != this.props.debug) {
            const term = this.refs.xterm.getTerminal()
            var socket = new WebSocket('ws://localhost:9000/pty?debug=' + this.props.debug);
            socket.onclose = (e) => {
              term.write("Connection closed: " + e.reason);
            }
            term.attach(socket)
            term.focus()
        }
    }
    render() {
        return <XTerm ref='xterm' />
    }

}

export default MyTerminal
