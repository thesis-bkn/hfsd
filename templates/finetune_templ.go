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
	"strings"

	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/templates/components"
)

type modelStatus int

const (
	Sampling modelStatus = iota
	Rating
	Finetuned
	Training
)

func (s modelStatus) ToolTip() string {
	switch s {
	case Sampling:
		return "This model is being sampling"
	case Rating:
		return "This model is sampled, please give feedback for finetuning"
	case Training:
		return "This model is being trained based on previous feedbacks"
	case Finetuned:
		return "This model is ready to be used to inference or further finetuning"
	}

	return ""
}

type ModelNode struct {
	ID     string
	Name   string
	Status modelStatus
	Parent *string
}

func FinetuneView(models []ModelNode, domain entity.Domain) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = treeStyle().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for i := range models {
				templ_7745c5c3_Err = swalStyle(&models[i], true, domain != entity.DomainSimpleAnimal, true).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"max-w-max mx-auto fixed inset-x-0 top-10 capitalize text-2xl text-gray-900 dark:text-white\"><p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(snakeToPascalCase(domain.String()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/finetune.templ`, Line: 49, Col: 42}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" Models</p></div><div class=\"flex overflow-x-auto justify-center mt-8\"><ul class=\"tree\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = tree(graph(models)).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></div>")
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

func snakeToPascalCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func swalFire(modelID string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_swalFire_33a2`,
		Function: `function __templ_swalFire_33a2(modelID){Swal.fire(
        { template: ` + "`" + `#model-${modelID}` + "`" + ` }
    ).then((result) => {
        if (result.isConfirmed) {
            fetch(` + "`" + `/api/finetune/${modelID}` + "`" + `, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
            }).then(response => {
                if (response.ok) {
                    Swal.fire({
                        icon: 'success',
                        title: 'Success!',
                        text: 'Submit finetuning model successfully',
                        confirmButtonColor: "#cccccc",
                        confirmButtonText: 'Continue'
                    }).then(_ => {
                        window.location.reload()
                    });
                } else {
                    Swal.fire({
                        icon: 'error',
                        title: 'Error!',
                        text: 'Failed to submit finetune new model.',
                        confirmButtonColor: "#cccccc",
                        confirmButtonText: 'Continue'
                    });
                }
            })

        } else if (result.isDenied) {
            window.location.replace(` + "`" + `/inference?model=${modelID}` + "`" + `)
        }
    })
}`,
		Call:       templ.SafeScript(`__templ_swalFire_33a2`, modelID),
		CallInline: templ.SafeScriptInline(`__templ_swalFire_33a2`, modelID),
	}
}

func redirectTo(url string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_redirectTo_047d`,
		Function: `function __templ_redirectTo_047d(url){window.location.replace(url)
}`,
		Call:       templ.SafeScript(`__templ_redirectTo_047d`, url),
		CallInline: templ.SafeScriptInline(`__templ_redirectTo_047d`, url),
	}
}

func tree(modelM map[string]*ModelNode, graphM map[string][]*ModelNode, curModel string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 = []any{"btn", "place-content-center", "duration-150", "ease-in-out", "hover:bg-primary-accent-200",
			templ.KV("btn-accent", modelM[curModel].Status == Sampling),
			templ.KV("btn-secondary", modelM[curModel].Status == Rating),
			templ.KV("btn-primary", modelM[curModel].Status == Finetuned),
			templ.KV("btn-warning", modelM[curModel].Status == Training)}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var5...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, swalFire(modelM[curModel].ID), redirectTo(fmt.Sprintf("/feedback/%s", modelM[curModel].ID)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var5).String()))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if modelM[curModel].Status == Finetuned {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" onClick=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 templ.ComponentScript = swalFire(modelM[curModel].ID)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var6.Call)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if modelM[curModel].Status == Rating {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" onClick=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 templ.ComponentScript = redirectTo(fmt.Sprintf("/feedback/%s", modelM[curModel].ID))
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var7.Call)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" data-twe-toggle=\"tooltip\" data-twe-placement=\"right\" data-twe-ripple-init data-twe-ripple-color=\"light\" title=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(modelM[curModel].Status.ToolTip()))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(modelM[curModel].Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/finetune.templ`, Line: 131, Col: 26}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(graphM[curModel]) > 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<ul>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, child := range graphM[curModel] {
				templ_7745c5c3_Err = tree(modelM, graphM, child.ID).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func graph(models []ModelNode) (map[string]*ModelNode, map[string][]*ModelNode, string) {
	graphM := make(map[string][]*ModelNode)
	modelM := make(map[string]*ModelNode)

	for i := range models {
		graphM[models[i].ID] = []*ModelNode{}
		modelM[models[i].ID] = &models[i]
	}

	var rootID string
	for i, model := range models {
		if model.Parent == nil || graphM[*model.Parent] == nil {
			rootID = model.ID
			continue
		}

		graphM[*model.Parent] = append(graphM[*model.Parent], &models[i])
	}

	return modelM, graphM, rootID
}

func swalStyle(model *ModelNode, shows ...bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var9 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var9 == nil {
			templ_7745c5c3_Var9 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<template id=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("model-%s", model.ID)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><swal-title><span class=\"text-blue-600\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(model.Name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/finetune.templ`, Line: 168, Col: 43}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></swal-title> <swal-html><p>Finetune or inference using this model</p></swal-html> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if shows[0] {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<swal-button type=\"confirm\" color=\"#36D399\">Finetune</swal-button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if shows[1] {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<swal-button type=\"deny\" color=\"#CC009C\">Inference</swal-button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if shows[2] {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<swal-button type=\"cancel\" color=\"#CCCCCC\">Cancel</swal-button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</template>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func treeStyle() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var11 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var11 == nil {
			templ_7745c5c3_Var11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<style type=\"text/css\">\n        .tree,\n        .tree ul,\n        .tree li {\n            list-style: none;\n            margin: 0;\n            padding: 0;\n            position: relative;\n        }\n        .tree {\n            margin: 0 0 1em;\n            text-align: center;\n        }\n        .tree,\n        .tree ul {\n            display: table;\n        }\n        .tree ul {\n            width: 100%;\n        }\n        .tree li {\n            display: table-cell;\n            padding: .5em 0;\n            vertical-align: top;\n        }\n        .tree li:before {\n            outline: solid 1px #666;\n            content: \"\";\n            left: 0;\n            position: absolute;\n            right: 0;\n            top: 0;\n        }\n        .tree li:first-child:before {\n            left: 50%;\n        }\n        .tree li:last-child:before {\n            right: 50%;\n        }\n        .tree code,\n        .tree button {\n            border: solid .1em #666;\n            border-radius: .2em;\n            display: inline-block;\n            margin: 0 .2em .5em;\n            padding: .2em .5em;\n            position: relative;\n        }\n        .tree ul:before,\n        .tree code:before,\n        .tree button:before {\n            outline: solid 1px #666;\n            content: \"\";\n            height: .5em;\n            left: 50%;\n            position: absolute;\n        }\n        .tree ul:before {\n            top: -.5em;\n        }\n        .tree code:before,\n        .tree button:before {\n            top: -.55em;\n        }\n        .tree>li {\n            margin-top: 0;\n        }\n        .tree>li:before,\n        .tree>li:after,\n        .tree>li>code:before,\n        .tree>li>button:before {\n            outline: none;\n        }\n    </style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
