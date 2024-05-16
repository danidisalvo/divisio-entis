import {Component, Inject} from '@angular/core';
import {FormBuilder, Validators} from '@angular/forms'
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import { v4 as uuidv4 } from 'uuid';
import {environment} from "../../environments/environment";
import {Node} from "../graph/graph.component";
import {HttpClient} from "@angular/common/http";

@Component({
  selector: 'app-edit-node-dialog',
  templateUrl: './edit-node-dialog.component.html',
  styleUrls: ['./edit-node-dialog.component.css']
})
export class EditNodeDialogComponent {

  parent?: string
  targetNodes?: Array<Node>

  form = this.fb.group({
    name: [
      {value: this.data.node.name, disabled: false},
      [Validators.required, Validators.pattern('(?! ).*[^ ]$'), Validators.maxLength(64)]],
    color: [
      {value: this.data.node.color, disabled: false},
      [Validators.required, Validators.pattern('^#(?:[0-9a-fA-F]{3}){1,2}$')]
    ],
    type: [{value: this.data.node.type, disabled: false}],
    targetNode: [{value: this.data.d.parent === null ? "" : this.data.d.parent.data.id, disabled: false}],
    addChild: [false],
    child: this.fb.group({
      name: ['', [Validators.pattern('(?! ).*[^ ]$'), Validators.maxLength(64)]],
      color: [
        {value: this.data.node.color, disabled: false},
        [Validators.required, Validators.pattern('^#(?:[0-9a-fA-F]{3}){1,2}$')]
      ],
      type: ['lexeme']
    })
  });

  constructor(private http: HttpClient,
              private fb: FormBuilder,
              public dialogRef: MatDialogRef<EditNodeDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any) {
    const url = `${environment.apiUrl}`;
    // this.parent = data.d.parent.data.id;
    this.http.get<Node[]>(url + "nodes/" + data.node.id + "/targets").subscribe({
      next: data => {
        this.targetNodes = data
      },
      error: error => {
        console.error(error.error.status + ': ' + error.error.message);
      }
    });
  }

  public onSubmit(): void {
    this.dialogRef.close({
      targetNode: this.form.controls['targetNode'].value,
      action: 'save',
      d: this.data.d,
      node: {
        id: this.data.node.id,
        name: this.form.controls['name'].value,
        color: this.form.controls['color'].value,
        type: this.form.controls['type'].value,
        properties: this.data.node.properties,
        children: [] // we leave them empty because the backend will not set them
      },
      child: {
        id: uuidv4(),
        name: this.form.controls['child'].controls['name'].value,
        color: this.form.controls['child'].controls['color'].value,
        type: this.form.controls['child'].controls['type'].value,
        properties: {},
        children: []
      },
      addChild: this.isChildToggled()
    });
  }

  public isChildToggled() {
    return this.form.controls['addChild'].value;
  }

  public toggleValidation() {
    this.form.controls['child'].controls['name'].setValidators(this.isChildToggled() ? [Validators.required] : null);
    this.form.controls['child'].controls['name'].updateValueAndValidity();

    this.form.controls['child'].controls['color'].setValidators(this.isChildToggled() ? [Validators.required] : null);
    this.form.controls['child'].controls['color'].updateValueAndValidity();
  }

  public isSavedEnabled() {
    return this.form.valid;
  }
}
