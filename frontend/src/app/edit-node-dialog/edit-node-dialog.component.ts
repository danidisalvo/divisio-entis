import {Component, Inject} from '@angular/core';
import {FormBuilder, Validators} from '@angular/forms'
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import { v4 as uuidv4 } from 'uuid';

@Component({
  selector: 'app-edit-node-dialog',
  templateUrl: './edit-node-dialog.component.html',
  styleUrls: ['./edit-node-dialog.component.css']
})
export class EditNodeDialogComponent {

  form = this.fb.group({
    color: [
      {value: this.data.node.color, disabled: false},
      [Validators.required, Validators.pattern('^#(?:[0-9a-fA-F]{3}){1,2}$')]
    ],
    addChild: [false],
    child: this.fb.group({
      name: ['', [Validators.pattern('(?! ).*[^ ]$'), Validators.maxLength(32)]],
      color: [
        {value: this.data.node.color, disabled: false},
        [Validators.required, Validators.pattern('^#(?:[0-9a-fA-F]{3}){1,2}$')]
      ],
      type: ['lexeme']
    })
  });

  constructor(private fb: FormBuilder,
              public dialogRef: MatDialogRef<EditNodeDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any) {
  }

  public onSubmit(): void {
    this.dialogRef.close({
      action: 'save',
      d: this.data.d,
      node: {
        id: this.data.node.id,
        name: this.data.node.name,
        color: this.form.controls['color'].value,
        type: this.data.node.type,
        properties: this.data.node.properties,
        children: []
      },
      child: {
        id: uuidv4(),
        name: this.form.controls['child'].controls['name'].value,
        color: this.form.controls['child'].controls['color'].value,
        type: this.form.controls['child'].controls['type'].value,
        properties: {},
        children: []
      },
      addChild: this.hasChild()
    });
  }

  public hasChild() {
    return this.form.controls['addChild'].value;
  }

  public toggleValidation() {
    this.form.controls['child'].controls['name'].setValidators(this.hasChild() ? [Validators.required] : null);
    this.form.controls['child'].controls['name'].updateValueAndValidity();

    this.form.controls['child'].controls['color'].setValidators(this.hasChild() ? [Validators.required] : null);
    this.form.controls['child'].controls['color'].updateValueAndValidity();
  }

  public isSavedEnabled() {
    return this.form.valid && (this.hasChild());
  }
}
