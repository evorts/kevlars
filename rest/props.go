/**
 * @Author: steven
 * @Description:
 * @File: props
 * @Version: 1.0.0
 * @Date: 07/06/23 20.16
 */

package rest

type props struct {
	headers map[string]string
	token   string
}

type Props interface {
	apply(m *props)
}

type propFunc func(*props)

func (fn propFunc) apply(p *props) {
	fn(p)
}

func ReplaceHeaders(v map[string]string) Props {
	return propFunc(func(p *props) {
		p.headers = v
	})
}

func AddHeader(k, v string) Props {
	return propFunc(func(p *props) {
		p.headers[k] = v
	})
}

func AppendToHeaders(values map[string]string) Props {
	return propFunc(func(p *props) {
		for k, v := range values {
			p.headers[k] = v
		}
	})
}

func newProps() *props {
	return &props{headers: map[string]string{}}
}
