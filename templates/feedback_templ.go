// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"

	"github.com/thesis-bkn/hfsd/templates/components"
)

type FeedbackAsset struct {
	ImageUrl string
	Group    int
	Order    int
}

func FeedBackView(modelID string, assets []FeedbackAsset) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			templ_7745c5c3_Err = components.NavBar("finetune").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <form id=\"feedbacks\"><input type=\"hidden\" id=\"modelID\" name=\"modelId\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(modelID))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <input type=\"hidden\" id=\"optionsLen\" name=\"optionsLen\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", len(assets))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"flex flex-wrap m-20 space-x-4 space-y-1 \">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, asset := range assets {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card w-96 bg-base-100 shadow-xl\"><figure><img src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
                fmt.Println("image url: ", asset.ImageUrl)
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(asset.ImageUrl))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" alt=\"Image\"></figure><div class=\"card-body\"><h2 class=\"card-title\"><div><label class=\"label cursor-pointer\"><input type=\"radio\" class=\"radio radio-accent\" id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("like-%d-%d", asset.Group, asset.Order)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" name=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("pref-%d-%d", asset.Group, asset.Order)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" checked value=\"like\"> <span class=\"label-text ml-1\">Like</span></label></div><div><label class=\"label cursor-pointer\"><input type=\"radio\" class=\"radio radio-error\" id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("dislike-%d-%d", asset.Group, asset.Order)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" name=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("pref-%d-%d", asset.Group, asset.Order)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" value=\"disklike\"> <span class=\"label-text ml-1\">Dislike</span></label></div></h2></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"max-w-max mx-auto fixed inset-x-0 bottom-10\"><button type=\"submit\" class=\"btn btn-active btn-secondary mt-5 px-20\">Submit</button></div></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = submitFeedbacks(modelID).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = base().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func submitFeedbacks(modelID string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n        // Handle form submission\n        document.getElementById('feedbacks').addEventListener('submit', function(event) {\n            event.preventDefault();\n            const formData = new FormData(this);\n            const modelID = document.getElementById('modelID').value; \n            const optionsLen = parseInt(document.getElementById('optionsLen').value);\n            const selectedItems = [];\n\n            for (const pair of formData.entries()) {\n                elements = pair[0].split('-');\n                if (elements.length != 3 || elements[0] != 'pref') {\n                    continue;\n                }\n\n                _group = elements[1]\n                order = elements[2]\n                option = pair[1] == 'like'\n\n                selectedItems.push({\n                    \"order\": parseInt(order),\n                    \"option\": option,\n                });\n            }\n\n            if (selectedItems.length != optionsLen) {\n                Swal.fire({\n                    title: 'Error!',\n                    text: 'Please select an option for each radio group',\n                    icon: 'error',\n                    confirmButtonText: 'OK'\n                });\n\n                return;\n            }\n\n            // Send selectedItems to server\n            fetch(`/api/feedback/${modelID}`, {\n                method: 'POST',\n                headers: {\n                    'Content-Type': 'application/json'\n                },\n                body: JSON.stringify({\n                    \"modelID\": modelID,\n                    \"items\": selectedItems\n                })\n            })\n            .then(response => {\n                if (response.ok) {\n                    Swal.fire({\n                        title: 'Success!',\n                        text: 'Success upload feedbacks',\n                        icon: 'success',\n                        confirmButtonColor: \"#cccccc\",\n                        confirmButtonText: 'Continue'\n                    }).then((result) => {\n                        window.location.replace('/factory')\n                    })\n                } else {\n                    Swal.fire({\n                        title: 'Error!',\n                        text: 'Failed to upload, please try again',\n                        icon: 'error',\n                        confirmButtonColor: \"#cccccc\",\n                        confirmButtonText: 'Continue'\n                    })\n                }\n            })\n            .catch(error => {\n                console.error('Error:', error);\n            });\n        });\n    </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
