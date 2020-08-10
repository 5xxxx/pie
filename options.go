/*
 *
 * options.go
 * tugrik
 *
 * Created by lintao on 2020/6/8 4:05 下午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package tugrik

type Options struct {
	UpdateEmpty bool
}

func (o *Options) SetUpdateEmpty(e bool) {
	o.UpdateEmpty = e
}
