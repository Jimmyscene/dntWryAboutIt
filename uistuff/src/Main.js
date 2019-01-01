import React, { Component, useState } from 'react';
import App from './App';
import './Main.css';
import TippyTappy from './TippyTappy'
import Code from './Code'
import { XTerm, Terminal } from 'react-xterm'

import * as  attach from './attach';
import * as fit from 'xterm/lib/addons/fit/fit';
import * as fullscreen from 'xterm/lib/addons/fullscreen/fullscreen';
import * as search from 'xterm/lib/addons/search/search';
import 'xterm/src/xterm.css'

Terminal.applyAddon(attach)
Terminal.applyAddon(fit)
Terminal.applyAddon(fullscreen)
Terminal.applyAddon(search)


class Main extends Component {
  constructor(props) {
    super(props)
    this.state = {
      'activeStep': 0,
      'debug': false,
      'activeSocket': null,
      'result': -1
      // 'initialized': false
    }
    this.setActiveStep = this.setActiveStep.bind(this)
    // this.setInitialized = this.setInitialized.bind(this)
    this.setDebug = this.setDebug.bind(this)
    this.runThings = this.runThings.bind(this)
  }

  setActiveStep(val) {
    this.setState((prevState, prevProps) =>  ({
      ...prevState, activeStep: val
    }))
  }
  setDebug(val) {
    this.setState((prevState, prevProps) =>  ({
      ...prevState, debug: val
    }))
  }
  runThings() {
    const term = this.refs.xterm.getTerminal()
    term.clear()

    if (this.state.activeSocket) {
      this.state.activeSocket.close()
    }
    var socket = new WebSocket('ws://localhost:9000/pty?debug=' + this.state.debug);
    this.setActiveStep(1)
    socket.onclose = (e) => {
      console.log(e)
      term.write("Connection closed: " + e.reason + "\r\n");

      this.setActiveStep(2)
    }
    term.attach(socket)
    term.focus()
    this.setState({
      activeSocket: socket
    })
  }
  render() {
    return (
      <div >
         <TippyTappy
          activeStep={this.state.activeStep}
        />
        <div className="row">
          <div className="column">
           <App result={this.state.result} debug={this.state.debug} debugHandler={this.setDebug} initailizedHandler={this.runThings}/>
          </div>
          <div className="column">
              <Code debug={this.state.debug} />
              <XTerm ref='xterm' />
          </div>
        </div>
      </div>
    );
  }

}

export default Main
