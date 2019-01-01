import React, { useState } from 'react';
import App from './App';
import './Main.css';
import TippyTappy from './TippyTappy'
import Code from './Code'
import MyTerminal from './terminal'

export default function Main() {
  const [activeStep, setActiveStep] = useState(0);
  const [debug, setdebug] = useState(false);
  const [initialized, setInitialized] = useState(false)
  return (
    <div >
       <TippyTappy
        activeStep={activeStep}
      />
      <div className="row">
        <div className="column">
         <App debug={debug} debugHandler={setdebug} initailizedHandler={setInitialized}/>
        </div>
        <div className="column">
            <Code debug={debug} />
            { initialized ? <MyTerminal debug={debug} /> : null }
        </div>
      </div>
    </div>
  );
}

