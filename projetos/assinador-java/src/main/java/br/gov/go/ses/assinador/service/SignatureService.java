package br.gov.go.ses.assinador.service;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.model.SignRequest;
import br.gov.go.ses.assinador.model.ValidateRequest;

public interface SignatureService {
    AssinadorResponse sign(SignRequest request);

    AssinadorResponse validate(ValidateRequest request);
}
