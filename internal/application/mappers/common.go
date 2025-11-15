package mappers

func EntitiesToDTOs[ENTITY any, DTO any](entities []*ENTITY, mapFunc func(*ENTITY) (*DTO, error)) ([]*DTO, error) {
	dtos := make([]*DTO, len(entities))
	for i, entity := range entities {
		dto, err := mapFunc(entity)
		if err != nil {
			return nil, err
		}
		dtos[i] = dto
	}

	return dtos, nil
}

func DTOsToEntities[ENTITY any, DTO any](dtos []*DTO, mapFunc func(*DTO) (*ENTITY, error)) ([]*ENTITY, error) {
	entities := make([]*ENTITY, len(dtos))
	for i, dto := range dtos {
		entity, err := mapFunc(dto)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
