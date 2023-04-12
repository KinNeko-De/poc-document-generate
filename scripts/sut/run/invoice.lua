local invoiceV1 = {
  deliveredOn = {year=2020, month=4, day=13}
}

local invoiceRecipientV1 = {
  name = "Max Mustermann",
  street = "Musterstraße 17",
  city = "Musterstadt",
  postCode = "12345",
  country = "DE"
}

local invoiceItemsV1 = {
  [1] = {
    description = "vfdsdsfdsfdsfs fdsfdskfdsk fdskfk fkwef kefkwekfe\\\\ANS 23054303053",
    quantity = 2,
    netAmount = 3.35,
    taxation = 19,
    totalAmount = 3.99,
    sum = 7.98,
    currency = "EUR"
  },
  [2] = {
    description = "vf ds dsf dsf dsfs fds fd skf dsk\\\\ANS 606406540",
    quantity = 1,
    netAmount = 9.07,
    taxation = 19,
    totalAmount = 10.79,
    sum = 10.79,
    currency = "EUR"
  },
  [3] = {
    description = "Versandkosten",
    quantity = 1,
    netAmount = 0.00,
    taxation = 0,
    totalAmount = 0.00,
    sum = 0.00,
    currency = "EUR"
  }
}

return { invoice = invoiceV1, invoiceRecipient = invoiceRecipientV1, invoiceItems = invoiceItemsV1 }





