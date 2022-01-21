import * as React from 'react';
import Chip from '@mui/material/Chip';


const getColor = (status) => {
    switch (status) {
        case 'failed':
            return 'error';
        case 'completed':
            return 'success';
        case 'scheduled':
            return 'info';
        default:
            return 'default';
    }
}

export default function ColoredStatus(props) {
  return (
    <Chip label={props.status} color={getColor(props.status)} />
  );
}