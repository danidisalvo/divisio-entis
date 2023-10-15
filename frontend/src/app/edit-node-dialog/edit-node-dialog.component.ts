import {Component, Inject} from '@angular/core';
import {FormBuilder, Validators} from '@angular/forms'
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";

@Component({
  selector: 'app-edit-node-dialog',
  templateUrl: './edit-node-dialog.component.html',
  styleUrls: ['./edit-node-dialog.component.css']
})
export class EditNodeDialogComponent {

  form = this.fb.group({
    color: [
      {value: this.data.node.color, disabled: this.data.node.readOnly},
      [Validators.required, Validators.pattern('^#(?:[0-9a-fA-F]{3}){1,2}$')]
    ],
    notes: [
      {value: this.data.node.notes, disabled: this.data.node.readOnly},
      [Validators.maxLength(512)]
    ],
    addChild: [false],
    child: this.fb.group({
      name: ['', [Validators.pattern('(?! ).*[^ ]$'), Validators.maxLength(32)]],
      color: ['', [Validators.pattern('^#(?:[0-9a-fA-F]{3}){1,2}$')]],
      notes: ['', [Validators.maxLength(512)]]
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
        name: this.data.node.name,
        color: this.form.controls['color'].value,
        notes: this.form.controls['notes'].value,
        children: []
      },
      child: {
        name: this.form.controls['child'].controls['name'].value,
        color: this.form.controls['child'].controls['color'].value,
        notes: this.form.controls['child'].controls['notes'].value,
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
    return this.form.valid && (this.hasChild() && this.data.node.readOnly || !this.data.node.readOnly);
  }
}
