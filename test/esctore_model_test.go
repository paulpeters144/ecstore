package ecstore_test

type SimpleEntity struct{ id string }

func (s *SimpleEntity) Id() string { return s.id }

type OtherEntity struct{ id string }

func (s *OtherEntity) Id() string { return s.id }

type ComplexEntity struct{ id string }

func (c *ComplexEntity) Id() string { return c.id }
