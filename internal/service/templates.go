package service

type MessageTemplate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type TemplateManager struct {
	templates map[string]MessageTemplate
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: map[string]MessageTemplate{
			"movement.created.entrada": {
				Title:   "✅ Ingreso de inventario registrado",
				Message: "Se registraron %d unidades del producto %s. Total en inventario: %d unidades.",
			},
			"movement.created.salida": {
				Title:   "📦 Salida de inventario registrada",
				Message: "Se retiraron %d unidades del producto %s. Quedan en inventario: %d unidades.",
			},
			"movement.updated.entrada": {
				Title:   "🔄 Ingreso actualizado",
				Message: "Se actualizó el ingreso a %d unidades del producto %s.",
			},
			"movement.updated.salida": {
				Title:   "🔄 Salida actualizada",
				Message: "Se actualizó la salida a %d unidades del producto %s.",
			},
		},
	}
}

func (tm *TemplateManager) Get(key string) (MessageTemplate, bool) {
	template, exists := tm.templates[key]
	return template, exists
}

func (tm *TemplateManager) GetAll() map[string]MessageTemplate {
	return tm.templates
}
