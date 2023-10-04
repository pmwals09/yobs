// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/pmwals09/yobs/internal/models/opportunity"
)

func OpportunityList(opportunities []opportunity.Opportunity) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<h2 class=\"text-[1.5em]\">")
		if err != nil {
			return err
		}
		var_2 := `Opportunities`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h2><table class=\"table-auto border-collapse w-full\" id=\"main-content\"><thead><tr><th class=\"border-b text-left\">")
		if err != nil {
			return err
		}
		var_3 := `Company Name`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
		if err != nil {
			return err
		}
		var_4 := `Role`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
		if err != nil {
			return err
		}
		var_5 := `URL`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
		if err != nil {
			return err
		}
		var_6 := `Status`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
		if err != nil {
			return err
		}
		var_7 := `Application Date`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
		if err != nil {
			return err
		}
		var_8 := `Actions`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</th></tr></thead><tbody>")
		if err != nil {
			return err
		}
		for _, o := range opportunities {
			_, err = templBuffer.WriteString("<tr class=\"hover:bg-sky-100\"><td class=\"border-b\">")
			if err != nil {
				return err
			}
			var var_9 string = o.CompanyName
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
			if err != nil {
				return err
			}
			var var_10 string = o.Role
			_, err = templBuffer.WriteString(templ.EscapeString(var_10))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
			if err != nil {
				return err
			}
			var var_11 string = o.URL
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
			if err != nil {
				return err
			}
			var var_12 string = string(o.Status)
			_, err = templBuffer.WriteString(templ.EscapeString(var_12))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
			if err != nil {
				return err
			}
			var var_13 string = formatApplicationDate(o.ApplicationDate)
			_, err = templBuffer.WriteString(templ.EscapeString(var_13))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</td><td class=\"border-b\"><a class=\"btn\" href=\"")
			if err != nil {
				return err
			}
			var var_14 templ.SafeURL = insertIDIntoHref("opportunities/{}", o.ID)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_14)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_15 := `Details`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a><!--")
			if err != nil {
				return err
			}
			var_16 := ` <a class="btn" href={ insertIDIntoHref("opportunities/{}/edit", o.ID) }>Edit</a> `
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("--></td></tr>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</tbody></table>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
