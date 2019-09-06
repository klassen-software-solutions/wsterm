//
// terminal/utility.go
// wsterm
//
// Created by steve on 2019-08-30.
// Copyright © 2019 Klassen Software Solutions. All rights reserved.
// Permission is hereby granted for use under the MIT License (https://opensource.org/licenses/MIT).
//

package terminal

import "io"

// CloseAndIgnore will close something ignoring any errors
func CloseAndIgnore(c io.Closer) {
	_ = c.Close()
}
