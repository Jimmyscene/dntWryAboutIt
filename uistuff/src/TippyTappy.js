import React from 'react';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';

const TippyTappy = (props) => {
    const steps = [
        "Select Test Type",
        "Run Test",
        "Test Finished"
    ]

    return (
        <Stepper activeStep={props.activeStep}>
            {steps.map((label, index) => {
              const stepProps = {};
              const labelProps = {};
              return (
                <Step key={label} {...stepProps}>
                  <StepLabel {...labelProps}>{label}</StepLabel>
                </Step>
              );
            })}
          </Stepper>
    )
}

export default TippyTappy
