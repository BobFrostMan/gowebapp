package model

import "errors"

// *****************************************************************************
//  Database communication errors
// *****************************************************************************

var NoDBConnection error = errors.New("Can't establish connection to database")
var DBNotSelected error = errors.New("Database type wasn't selected")


